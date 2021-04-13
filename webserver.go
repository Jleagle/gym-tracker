package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/Jleagle/pure-gym-tracker/helpers"
	"github.com/Jleagle/pure-gym-tracker/influx"
	influxquerybuilder "github.com/benjamin658/influx-query-builder"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"go.uber.org/zap"
)

func webserver() {

	app := fiber.New()

	// Middleware
	app.Use(cache.New(cache.Config{Expiration: time.Minute, CacheControl: true}))
	app.Use(compress.New(compress.Config{Level: compress.LevelBestSpeed}))

	// Routes
	app.Get("/", rootHandler)
	app.Get("/heatmap.json", heatmapHandler)
	app.Post("/submit", submitHandler)

	// Serve
	err := app.Listen("0.0.0.0:" + os.Getenv("PURE_PORT_BACKEND"))
	if err != nil {
		logger.Error("serving webserver", zap.Error(err))
	}
}

func rootHandler(c *fiber.Ctx) error {
	return c.SendString("OK")
}

func heatmapHandler(c *fiber.Ctx) error {

	q := influxquerybuilder.New().
		Select(`max("max") AS max_max`).
		Select(`max("pcnt") AS max_pcnt`).
		Select(`max("people") AS max_people`).
		From(`PureGym"."alltime"."gyms`).
		Where(`time`, `>`, `now()-1d`).
		Where(`gym`, `=`, `Fareham`).
		GroupByTime(influxquerybuilder.NewDuration().Minute(10)).
		Fill("null")

	resp, err := influx.Read(q.Build())
	if err != nil {
		logger.Error("failed to query influx", zap.Error(err))
	}

	var hc = map[string][][]interface{}{}

	var data = map[time.Weekday]map[int][]float64{}

	if len(resp.Results) > 0 && len(resp.Results[0].Series) > 0 {

		var series = resp.Results[0].Series[0]

		for k := range series.Columns {
			if k > 0 {

				for _, vv := range series.Values {

					t, err := time.Parse(time.RFC3339, vv[0].(string))
					if err != nil {
						logger.Error("casting", zap.Error(err))
						continue
					}

					val, err := vv[k].(json.Number).Float64()
					if err != nil {
						logger.Error("casting", zap.Error(err))
						continue
					}

					if data[t.Weekday()] == nil {
						data[t.Weekday()] = map[int][]float64{}
					}

					data[t.Weekday()][t.Hour()] = append(data[t.Weekday()][t.Hour()], val)
				}
			}
		}
	}

	for day, hours := range data {
		for hour, vals := range hours {
			hc["max_pcnt"] = append(hc["max_pcnt"], []interface{}{hour, day, helpers.Max(vals...)})
		}
	}

	return c.JSON(hc)
}

func submitHandler(c *fiber.Ctx) error {
	return c.SendString("new gym")
}
