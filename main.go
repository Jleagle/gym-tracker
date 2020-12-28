package main

import (
	"github.com/chromedp/cdproto/network"
	"github.com/robfig/cron/v3"
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
