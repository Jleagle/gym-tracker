package main

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Jleagle/gym-tracker/config"
	"github.com/Jleagle/gym-tracker/influx"
	"github.com/cenkalti/backoff/v4"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"go.uber.org/zap"
)

var (
	baseContext  context.Context
	membersRegex = regexp.MustCompile(`(?i)([0-9,]{1,4})\s(of)\s([0-9,]{1,4})`)
	cookies      []*network.Cookie
)

func init() {

	// create a new browser
	baseContext, _ = chromedp.NewContext(context.Background())

	// start the browser without a timeout
	err := chromedp.Run(baseContext)
	if err != nil {
		logger.Error("failed to start browser", zap.Error(err))
	}
}

func trigger() {

	peopleString, town, err, _ := loginAndCheckMembers(config.User, config.Pass)
	if err != nil {
		logger.Error("running chromedp", zap.Error(err))
		return
	}

	if town == "" {
		logger.Error("missing town")
		return
	}

	if peopleString == "10 or fewer people" {
		logger.Info("members", zap.String("town", town), zap.Int("now", 0))
		return
	}

	members := membersRegex.FindStringSubmatch(peopleString)
	if len(members) != 4 {
		logger.Error("parsing count failed", zap.String("string", peopleString))
		return
	}

	now, err := strconv.Atoi(strings.Replace(members[1], ",", "", 1))
	if err != nil {
		logger.Error("parsing members", zap.Error(err), zap.String("string", peopleString))
		return
	}

	max, err := strconv.Atoi(strings.Replace(members[3], ",", "", 1))
	if err != nil {
		logger.Error("parsing members", zap.Error(err), zap.String("string", peopleString))
		return
	}

	pct := calculatePercent(now, max)

	logger.Info("members", zap.String("town", town), zap.Int("max", max), zap.Int("now", now), zap.Float64("pct", pct))

	_, err = influx.Write(town, now, max, pct, time.Now())
	if err != nil {
		logger.Error("sending to influx failed", zap.Error(err))
	}
}

func loginAndCheckMembers(user, pass string) (people, gym string, err error, errorString string) {

	actions := []chromedp.Action{
		network.Enable(),
		chromedp.Emulate(device.IPadPro),
		chromedp.Navigate("https://www.puregym.com/members/"),
		chromedp.WaitVisible("input[name=username], input[name=password], #people_in_gym"),
		chromedp.ActionFunc(func(ctx context.Context) error {

			// Accept cookies, probably don't need to bother
			var cookieNodes []*cdp.Node
			err = chromedp.Nodes("button.coi-banner__accept", &cookieNodes, chromedp.AtLeast(0)).Do(ctx)
			if err != nil {
				return err
			}

			if len(cookieNodes) > 0 {

				// logger.Info("Submitting cookie popup")

				_, exp, err := runtime.Evaluate("CookieInformation.submitAllCategories();").Do(ctx)
				if err != nil {
					return err
				}
				if exp != nil {
					return exp
				}
			}

			// Login
			var loginNodes []*cdp.Node
			err = chromedp.Nodes("input[name=username], input[name=password]", &loginNodes, chromedp.AtLeast(0)).Do(ctx)
			if err != nil {
				return err
			}

			if len(loginNodes) > 0 {

				logger.Info("Logging in")

				err = chromedp.SendKeys("input[name=username]", user).Do(ctx)
				if err != nil {
					return err
				}

				err = chromedp.SendKeys("input[name=password]", pass).Do(ctx)
				if err != nil {
					return err
				}

				err = chromedp.Click("button[value=login]", chromedp.ByQuery).Do(ctx)
				if err != nil {
					return err
				}
			}

			return nil
		}),
		chromedp.WaitVisible("#people_in_gym, div.danger"),
		chromedp.ActionFunc(func(ctx context.Context) error {

			// Find error message on failure
			var errorNodes []*cdp.Node
			err = chromedp.Nodes("div.danger ul li", &errorNodes, chromedp.AtLeast(0)).Do(ctx)
			if err != nil {
				return err
			}

			if len(errorNodes) > 0 {

				err = chromedp.InnerHTML("div.danger ul li", &errorString).Do(ctx)
				if err != nil {
					return err
				}
			}

			// Save count on success
			var peopleNodes []*cdp.Node
			err = chromedp.Nodes("#people_in_gym", &peopleNodes, chromedp.AtLeast(0)).Do(ctx)
			if err != nil {
				return err
			}

			if len(peopleNodes) > 0 {

				err = chromedp.InnerHTML("#people_in_gym span", &people).Do(ctx)
				if err != nil {
					return err
				}

				err = chromedp.AttributeValue("#people_in_gym a", "href", &gym, nil).Do(ctx)
				if err != nil {
					return err
				}
			}

			//
			return nil
		}),
	}

	// Make context
	ctx, cancel1 := context.WithTimeout(baseContext, 30*time.Second)
	defer cancel1()

	ex, err := os.Executable()
	if err != nil {
		logger.Error("failed to find exe path", zap.Error(err))
		return
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.72 Safari/537.36"),
		chromedp.UserDataDir(filepath.Dir(ex)+"/user-data/"+credential.Email),
		chromedp.WindowSize(1920, 1080),
	)

	ctx, cancel2 := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel2()

	ctx, cancel3 := chromedp.NewContext(ctx)
	defer cancel3()

	// Retry
	work := func() error {
		return chromedp.Run(ctx, actions...)
	}

	scrape := backoff.Retry(work, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 10))

	return people, path.Base(gym), scrape, errorString
}
