package main

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

	database  = os.Getenv("PURE_INFLUX_DATABASE")
	retention = os.Getenv("PURE_INFLUX_RETENTION")
)

func getInfluxClient() (*influx.Client, error) {

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

func influxWriteMany(batch influx.BatchPoints) (resp *influx.Response, err error) {

	if len(batch.Points) == 0 {
		return nil, nil
	}

	batch.Database = database
	batch.RetentionPolicy = retention
	batch.Precision = batch.Points[0].Precision // Must be in batch and point

	if batch.Time.IsZero() || batch.Time.Unix() == 0 {
		batch.Time = time.Now()
	}

	client, err := getInfluxClient()
	if err != nil {
		return nil, err
	}

	return client.Write(batch)
}

func influxQuery(query string) (resp *influx.Response, err error) {

	client, err := getInfluxClient()
	if err != nil {
		return resp, err
	}

	resp, err = client.Query(influx.Query{
		Command:         query,
		Database:        database,
		RetentionPolicy: retention,
	})

	return resp, err
}
