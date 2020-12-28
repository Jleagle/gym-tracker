package main

import (
	"context"
	"os"
	"regexp"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var (
	regexMembers = regexp.MustCompile(`(?i)([0-9]{1,3}) of ([0-9]{1,3})`)
	cookies      []*network.Cookie
	logger       *zap.Logger
)

func init() {
	logger, _ = zap.NewDevelopment()
}

func main() {

	defer logger.Sync()

	c := cron.New()
	_, err := c.AddFunc("@every 10m", trigger)
	if err != nil {
		logger.Error("adding cron", zap.Error(err))
		return
	}
	c.Start()

	<-make(chan struct{})
}

func trigger() {

	ctx := context.Background()

	ctx, cancel1 := chromedp.NewContext(ctx)
	defer cancel1()

	ctx, cancel2 := context.WithTimeout(ctx, 30*time.Second)
	defer cancel2()

	people, town, err := loginAndCheckMembers(ctx)
	if err != nil {
		logger.Error("running chromedp", zap.Error(err))
		return
	}

	logger.Info("people ", people)
	logger.Info("town ", town)
}

func loginAndCheckMembers(ctx context.Context) (people, town string, err error) {

	actions := []chromedp.Action{
		network.Enable(),
		chromedp.ActionFunc(func(ctx context.Context) error {
			if len(cookies) > 0 {
				logger.Info("Setting cookies")
			}
			for _, v := range cookies {
				expr := cdp.TimeSinceEpoch(time.Unix(int64(v.Expires), 0))
				_, err := network.SetCookie(v.Name, v.Value).
					WithExpires(&expr).
					WithDomain(v.Domain).
					WithHTTPOnly(v.HTTPOnly).
					WithPath(v.Path).
					WithPriority(v.Priority).
					WithSameSite(v.SameSite).
					WithSecure(v.Secure).
					Do(ctx)

				if err != nil {
					return err
				}
			}
			return nil
		}),
		chromedp.Emulate(device.IPadPro),
		chromedp.Navigate("https://www.puregym.com/members/"),
		chromedp.WaitVisible("#loginForm, #people_in_gym"),
		chromedp.ActionFunc(func(ctx context.Context) error {

			// Accept cookies
			var cookieNodes []*cdp.Node
			err = chromedp.Nodes("button.coi-banner__accept", &cookieNodes, chromedp.AtLeast(0)).Do(ctx)
			if err != nil {
				return err
			}

			if len(cookieNodes) > 0 {

				logger.Info("Submitting cookie popup")
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
			err = chromedp.Nodes("#loginForm", &loginNodes, chromedp.AtLeast(0)).Do(ctx)
			if err != nil {
				return err
			}

			if len(loginNodes) > 0 {

				logger.Info("Logging in")

				err = chromedp.SendKeys("#loginForm input[type=email]", os.Getenv("PURE_USER")).Do(ctx)
				if err != nil {
					return err
				}

				err = chromedp.SendKeys("#loginForm input[type=password]", os.Getenv("PURE_PASS")).Do(ctx)
				if err != nil {
					return err
				}

				err = chromedp.Click("#login-submit", chromedp.ByID).Do(ctx)
				if err != nil {
					return err
				}
			}

			return nil
		}),
		chromedp.WaitVisible("#people_in_gym"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("Logged in, taking cookies")
			var err error
			cookies, err = network.GetAllCookies().Do(ctx)
			return err
		}),
		chromedp.InnerHTML("#people_in_gym span", &people),
		chromedp.InnerHTML("#people_in_gym a", &town),
	}

	err = chromedp.Run(ctx, actions...)
	return people, town, err
}
