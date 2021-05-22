package main

import (
	"context"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Jleagle/gym-tracker/datastore"
	"github.com/Jleagle/gym-tracker/influx"
	"github.com/Jleagle/gym-tracker/log"
	"github.com/cenkalti/backoff/v4"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"go.uber.org/zap"
)

var browserCtx context.Context

func init() {

	allocatorCtx, _ := chromedp.NewExecAllocator(context.Background(), chromedp.DefaultExecAllocatorOptions[:]...)

	browserCtx, _ = chromedp.NewContext(allocatorCtx)

	// Start a browser
	if err := chromedp.Run(browserCtx); err != nil {
		log.Instance.Error("starting browser", zap.Error(err))
	}
}

func scrapeGyms() {

	creds, err := datastore.GetCredentials()
	if err != nil {
		log.Instance.Error("failed to start browser", zap.Error(err))
		return
	}

	for _, v := range creds {
		scrapeGym(v)
	}
}

var (
	membersRegex = regexp.MustCompile(`(?i)([0-9,]{1,4})\s(of)\s([0-9,]{1,4})`)
	cookies      = map[string][]*network.Cookie{}
)

func scrapeGym(credential datastore.Credential) {

	peopleString, town, err, _ := scrape(credential)
	if err != nil {
		log.Instance.Error("running chromedp", zap.Error(err))
		return
	}

	if town == "" {
		log.Instance.Error("missing town")
		return
	}

	if peopleString == "10 or fewer people" {

		_, err = influx.Write(town, 0, 0, 0, time.Now())
		if err != nil {
			log.Instance.Error("sending to influx failed", zap.Error(err))
		}

		return
	}

	members := membersRegex.FindStringSubmatch(peopleString)
	if len(members) != 4 {
		log.Instance.Error("parsing count failed", zap.String("string", peopleString))
		return
	}

	now, err := strconv.Atoi(strings.Replace(members[1], ",", "", 1))
	if err != nil {
		log.Instance.Error("parsing members", zap.Error(err), zap.String("string", peopleString))
		return
	}

	max, err := strconv.Atoi(strings.Replace(members[3], ",", "", 1))
	if err != nil {
		log.Instance.Error("parsing members", zap.Error(err), zap.String("string", peopleString))
		return
	}

	pct := calculatePercent(now, max)

	_, err = influx.Write(town, now, max, pct, time.Now())
	if err != nil {
		log.Instance.Error("sending to influx failed", zap.Error(err))
	}
}

func scrape(credential datastore.Credential) (people, gym string, err error, errorString string) {

	actions := []chromedp.Action{
		network.Enable(),
		chromedp.ActionFunc(func(ctx context.Context) error {

			// logger.Info("Setting cookies", zap.Int("count", len(cookies)))

			network.ClearBrowserCookies()

			for _, cookie := range cookies[credential.Email] {

				expr := cdp.TimeSinceEpoch(time.Unix(int64(cookie.Expires), 0))
				err := network.SetCookie(cookie.Name, cookie.Value).
					WithExpires(&expr).
					WithDomain(cookie.Domain).
					WithHTTPOnly(cookie.HTTPOnly).
					WithPath(cookie.Path).
					WithPriority(cookie.Priority).
					WithSameSite(cookie.SameSite).
					WithSecure(cookie.Secure).
					WithSameParty(cookie.SameParty).
					WithSourcePort(cookie.SourcePort).
					WithSourceScheme(cookie.SourceScheme).
					Do(ctx)

				if err != nil {
					return err
				}
			}

			return nil
		}),
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

				// log.Instance.Info("Submitting cookie popup")

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

				log.Instance.Info("Logging in")

				err = chromedp.SendKeys("input[name=username]", credential.Email).Do(ctx)
				if err != nil {
					return err
				}

				err = chromedp.SendKeys("input[name=password]", credential.PIN).Do(ctx)
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
		chromedp.ActionFunc(func(ctx context.Context) error {

			// logger.Info("Logged in, taking cookies")

			var err error
			cookies[credential.Email], err = network.GetAllCookies().Do(ctx)
			return err
		}),
	}

	// Retry
	work := func() error {

		timeoutCtx, timeoutCancel := context.WithTimeout(browserCtx, 30*time.Second)
		defer timeoutCancel()

		return chromedp.Run(timeoutCtx, actions...)
	}

	notify := func(error, time.Duration) { log.Instance.Info("Failed to request counts", zap.Error(err)) }
	err = backoff.RetryNotify(work, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 10), notify)

	return people, path.Base(gym), err, errorString
}
