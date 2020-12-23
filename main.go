package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/chromedp/chromedp"
)

const (
	port        = 9222
	username    = ""
	password    = ""
	loginPath   = "https://www.puregym.com/Login/"
	membersPath = "https://www.puregym.com/members/"
)

var regexMembers = regexp.MustCompile(`([0-9]{1,3}) (of|OF) ([0-9]{1,3})`)

func main() {

	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var people string
	var town string

	actions := []chromedp.Action{
		chromedp.Navigate(membersPath),
		// chromedp.WaitVisible("#people_in_gym"),
		chromedp.WaitReady("#people_in_gym"),
		chromedp.InnerHTML("#people_in_gym span", &people),
		chromedp.InnerHTML("#people_in_gym a", &town),
	}

	err := chromedp.Run(ctx, actions...)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(people)
	fmt.Println(town)
}
