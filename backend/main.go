package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/Jleagle/gym-tracker/config"
	"github.com/Jleagle/gym-tracker/log"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func main() {

	rand.Seed(time.Now().Unix())

	defer func() {
		err := log.Instance.Sync()
		if err != nil {
			fmt.Println(err)
		}
	}()

	// Set timezone to UK
	loc, err := time.LoadLocation("Europe/London")
	if err != nil {
		log.Instance.Error("setting timezone", zap.Error(err))
	}
	time.Local = loc

	//importFromChronograf()
	//return

	if config.PortBackend == "" ||
		config.InfluxURL == "" ||
		config.InfluxUser == "" ||
		config.InfluxPass == "" ||
		config.GoogleProject == "" {
		log.Instance.Error("missing configs")
		return
	}

	// Flags
	disableScraping := flag.Bool("noscrape", false, "Disable scraping")
	flag.Parse()

	// Scrape
	if !*disableScraping {

		scrapeGyms()

		c := cron.New(cron.WithSeconds())
		_, err := c.AddFunc("10 */10 * * * *", scrapeGyms)
		if err != nil {
			log.Instance.Error("adding cron", zap.Error(err))
			return
		}
		c.Start()
	}

	// Serve JSON
	err = webserver()
	if err != nil {
		log.Instance.Error("serving webserver", zap.Error(err))
	}

	// Block
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func calculatePercent(now int, max int) float64 {
	return math.Round((float64(now)/float64(max))*100*100) / 100
}
