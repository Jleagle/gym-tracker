package influx

import (
	"net/url"
	"os"
	"sync"
	"time"

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

		host, err = url.Parse(os.Getenv("PURE_INFLUX_URL"))
		if err != nil {
			return nil, err
		}

		client, err = influx.NewClient(influx.Config{
			URL:      *host,
			Username: os.Getenv("PURE_INFLUX_USER"),
			Password: os.Getenv("PURE_INFLUX_PASS"),
		})
	}

	return client, err
}

func Write(gym string, count int, max int) (resp *influx.Response, err error) {

	batch := influx.BatchPoints{
		Points: []influx.Point{{
			Measurement: "gyms",
			Tags:        map[string]string{"gym": gym},
			Time:        time.Time{},
			Precision:   "m",
			Fields: map[string]interface{}{
				"people": count,
				"max":    max,
				"pcnt":   (float64(count) / float64(max)) * 100,
			},
		}},
		Database:        os.Getenv("PURE_INFLUX_DATABASE"),
		RetentionPolicy: os.Getenv("PURE_INFLUX_RETENTION"),
		Precision:       "m",
		Time:            time.Now(),
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
		Database:        os.Getenv("PURE_INFLUX_DATABASE"),
		RetentionPolicy: os.Getenv("PURE_INFLUX_RETENTION"),
	})

	return resp, err
}
