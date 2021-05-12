package influx

import (
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/Jleagle/gym-tracker/config"
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

func Write(gym string, count int, max int, percent float64, t time.Time) (resp *influx.Response, err error) {

	client, err := getClient()
	if err != nil {
		return nil, err
	}

	_, week := t.ISOWeek()

	return client.Write(influx.BatchPoints{
		Points: []influx.Point{{
			Measurement: "gyms",
			Tags: map[string]string{
				"gym":      gym,
				"yearDay":  strconv.Itoa(t.YearDay()),
				"yearWeek": strconv.Itoa(week),
				"monthDay": strconv.Itoa(t.Day()),
				"weekDay":  strconv.Itoa(int(t.Weekday())),
				"weekHour": strconv.Itoa(int(t.Weekday())) + "-" + strconv.Itoa(t.Hour()),
				"dayHour":  strconv.Itoa(t.Hour()),
			},
			Time:      t,
			Precision: "m",
			Fields: map[string]interface{}{
				"people": count,
				"max":    max,
				"pcnt":   percent,
			},
		}},
		Database:        config.InfluxDatabase,
		RetentionPolicy: config.InfluxRetention,
		Precision:       "m",
		Time:            t,
	})
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
