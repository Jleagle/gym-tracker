package main

import (
	"fmt"
	"regexp"

	"github.com/chromedp/cdproto/network"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var (
	membersRegex = regexp.MustCompile(`(?i)([0-9,]{1,4})\s(of)?\s?([0-9,]{1,4})?`)
	cookies      []*network.Cookie
	logger       *zap.Logger
)

func main() {

	if config.PortBackend == "" ||
		config.User == "" ||
		config.Pass == "" ||
		config.InfluxURL == "" ||
		config.InfluxUser == "" ||
		config.InfluxPass == "" ||
		config.InfluxDatabase == "" ||
		config.InfluxRetention == "" {
		logger.Error("missing configs")
		return
	}

	// Logger
	logger, _ = zap.NewDevelopment()

	defer func() {
		err := logger.Sync()
		if err != nil {
			fmt.Println(err)
		}
	}()

	// Update
	trigger()

	c := cron.New()
	_, err := c.AddFunc("*/10 * * * *", trigger)
	if err != nil {
		logger.Error("adding cron", zap.Error(err))
		return
	}
	c.Start()

	// Serve
	webserver()
}
