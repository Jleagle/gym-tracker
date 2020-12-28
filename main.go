package main

import (
	"fmt"
	"github.com/chromedp/cdproto/network"
	"go.uber.org/zap"
	"regexp"
)

var (
	membersRegex = regexp.MustCompile(`(?i)([0-9,]{1,4}) of ([0-9,]{1,4})`)
	cookies      []*network.Cookie
	logger       *zap.Logger
)

func init() {
	logger, _ = zap.NewDevelopment()
}

func main() {

	defer func() {
		err := logger.Sync()
		if err != nil {
			fmt.Println(err)
		}
	}()

	startCron()
	webserver()
}
