package influx

import (
	"math"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/Jleagle/puregym-tracker/config"
	influx "github.com/influxdata/influxdb1-client"
)

var (
	client *influx.Client
	mutex  sync.Mutex
)

func getClient() (*influx.Client, error) {

	mutex.Lock()
	defer mutex.Unlock()

	var err error
	var host *url.URL

	if client == nil {

		host, err = url.Parse(config.InfluxURL)
		if err != nil {
			return nil, err
		}

		client, err = influx.NewClient(influx.Config{
			URL:      *host,
			Username: config.InfluxUser,
			Password: config.InfluxPass,
		})
	}

	return client, err
}

func Write(gym string, count int, max int) (resp *influx.Response, err error) {

	t := time.Now()

	batch := influx.BatchPoints{
		Points: []influx.Point{{
			Measurement: "gyms",
			Tags: map[string]string{
				"gym":      gym,
				"yearDay":  strconv.Itoa(t.YearDay()),
				"monthDay": strconv.Itoa(t.Day()),
				"weekDay":  strconv.Itoa(int(t.Weekday())),
				"weekHour": strconv.Itoa(int(t.Weekday())) + "-" + strconv.Itoa(t.Hour()),
			},
			Time:      t,
			Precision: "m",
			Fields: map[string]interface{}{
				"people": count,
				"max":    max,
				"pcnt":   math.Round((float64(count)/float64(max))*100*100) / 100, // Percent, rounded
			},
		}},
		Database:        config.InfluxDatabase,
		RetentionPolicy: config.InfluxRetention,
		Precision:       "m",
		Time:            t,
	}

	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.Write(batch)
}

func Read(query string) (resp *influx.Response, err error) {

	client, err := getClient()
	if err != nil {
		return resp, err
	}

	resp, err = client.Query(influx.Query{
		Command:         query,
		Database:        config.InfluxDatabase,
		RetentionPolicy: config.InfluxRetention,
	})

	return resp, err
}
