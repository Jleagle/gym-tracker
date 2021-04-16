package main

import (
	"fmt"

	"github.com/Jleagle/puregym-tracker/config"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
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

	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc("10 */10 * * * *", trigger)
	if err != nil {
		logger.Error("adding cron", zap.Error(err))
		return
	}
	c.Start()

	// Serve
	webserver()
}
