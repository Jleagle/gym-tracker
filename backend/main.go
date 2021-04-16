package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/Jleagle/puregym-tracker/config"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func main() {

	// Logger
	logger, _ = zap.NewDevelopment()

	defer func() {
		err := logger.Sync()
		if err != nil {
			fmt.Println(err)
		}
	}()

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

	disableScraping := flag.Bool("noscrape", false, "Disable scraping")
	flag.Parse()

	// Scrape
	if !*disableScraping {

		trigger()

		c := cron.New(cron.WithSeconds())
		_, err := c.AddFunc("10 */10 * * * *", trigger)
		if err != nil {
			logger.Error("adding cron", zap.Error(err))
			return
		}
		c.Start()
	}

	// Serve JSON
	err := webserver()
	if err != nil {
		logger.Error("serving webserver", zap.Error(err))
	}

	// Block
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
