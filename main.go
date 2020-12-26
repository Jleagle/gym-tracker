package main

import (
	"context"
	"github.com/chromedp/cdproto/runtime"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"go.uber.org/zap"
)

const (
	username    = ""
	password    = ""
	membersPath = "https://www.puregym.com/members/"
)

var regexMembers = regexp.MustCompile(`([0-9]{1,3}) (of|OF) ([0-9]{1,3})`)
var cookies []*network.Cookie
var logger *zap.SugaredLogger

func init() {

	loggerDev, _ := zap.NewDevelopment()
	logger = loggerDev.Sugar()
}

func main() {

	defer logger.Sync()

	ctx := context.Background()

	ctx, cancel1 := chromedp.NewContext(ctx)
	defer cancel1()

	ctx, cancel2 := context.WithTimeout(ctx, 5*time.Second)
	defer cancel2()

	listenForNetworkEvent(ctx, logger)

	// people, town, err := checkMembers(ctx)
	// if errors.Is(err, context.DeadlineExceeded) {
	// 	logger.Error(err)
	people, town, err := loginAndCheckMembers(ctx)
	if err != nil {
		logger.Error(err)
	}
	// }

	if err := ioutil.WriteFile("screenshot1.png", b1, 0644); err != nil {
		log.Fatal(err)
	}

	logger.Info("people ", people)
	logger.Info("town ", town)
}

func checkMembers(ctx context.Context) (people, town string, err error) {

	actions := []chromedp.Action{
		network.Enable(),
		chromedp.Emulate(device.IPadPro),
		chromedp.ActionFunc(setCookies),
		chromedp.Navigate(membersPath),
		chromedp.WaitVisible("#people_in_gym"),
		chromedp.InnerHTML("#people_in_gym span", &people),
		chromedp.InnerHTML("#people_in_gym a", &town),
	}

	err = chromedp.Run(ctx, actions...)
	return people, town, err
}

var b1 []byte

func loginAndCheckMembers(ctx context.Context) (people, town string, err error) {

	actions := []chromedp.Action{
		network.Enable(),
		chromedp.Emulate(device.IPadPro),
		chromedp.Navigate(membersPath),
		chromedp.ActionFunc(func(ctx context.Context) error {

			_, exp, err := runtime.Evaluate("CookieInformation.submitAllCategories();").Do(ctx)
			if err != nil {
				return err
			}
			if exp != nil {
				return exp
			}
			return nil
		}),
		chromedp.WaitVisible("#loginForm"),
		chromedp.SendKeys("#loginForm input[type=email]", ""), // todo
		chromedp.SendKeys("#loginForm input[type=password]", ""), // todo
		chromedp.Submit("#loginForm"),
		chromedp.CaptureScreenshot(&b1),

		// CookieInformation.submitAllCategories();

		// chromedp.Click("#loginForm button[type=submit]"),
		// chromedp.Sleep(time.Second * 2),
		// chromedp.WaitVisible("#people_in_gym"),
		// chromedp.CaptureScreenshot(&b1),
		// chromedp.InnerHTML("#people_in_gym span", &people),
		// chromedp.InnerHTML("#people_in_gym a", &town),
		// chromedp.ActionFunc(getCookies),
	}

	err = chromedp.Run(ctx, actions...)
	return people, town, err
}

func getCookies(ctx context.Context) error {

	var err error
	cookies, err = network.GetAllCookies().Do(ctx)
	for _, v := range cookies {
		logger.Info(v)
	}

	return err
}

func setCookies(ctx context.Context) error {

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
			logger.Error(err)
		}
	}

	return nil
}

// var r1 = regexp.MustCompile(`(?i)/(Login|members)/`)

func listenForNetworkEvent(ctx context.Context, logger *zap.SugaredLogger) {
	chromedp.ListenTarget(ctx, func(ev interface{}) {

		switch ev := ev.(type) {
		case *network.EventResponseReceived:

			// if r1.MatchString(ev.Response.URL) {
			if strings.Contains(ev.Response.URL, "www.puregym.com/") {

				// logger.Info("network: ", ev.Response.URL)

				// val, ok := ev.Response.Headers["set-cookie"]
				// logger.Info(ok, val)

			}
		}
	})
}
