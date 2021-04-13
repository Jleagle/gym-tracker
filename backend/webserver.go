package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/Jleagle/pure-gym-tracker/influx"
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
	app.Get("/people.json", peopleHandler)
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

func peopleHandler(c *fiber.Ctx) error {

	var query = `SELECT mean("people") AS "mean_people", mean("pcnt") AS "mean_pcnt" FROM "PureGym"."alltime"."gyms" `

	switch c.Query("range") {
	case "year":
		query += "WHERE time > now()-365d GROUP BY yearDay"
	case "month":
		query += "WHERE time > now()-365d GROUP BY monthDay"
	case "week":
		query += "WHERE time > now()-365d GROUP BY weekDay"
	default:
		return nil
	}

	query += " FILL(0)"

	resp, err := influx.Read(query)
	if err != nil {
		logger.Error("failed to query influx", zap.Error(err))
	}

	var ret = map[string]map[string][]json.Number{}

	for _, result := range resp.Results {
		for _, series := range result.Series {
			for _, tagValue := range series.Tags {
				for kk, column := range series.Columns {
					if kk > 0 {

						if ret[tagValue] == nil {
							ret[tagValue] = map[string][]json.Number{}
						}
						if ret[tagValue][column] == nil {
							ret[tagValue][column] = []json.Number{}
						}
						ret[tagValue][column] = append(ret[tagValue][column], series.Values[0][kk].(json.Number))
					}
				}
			}
		}
	}

	return c.JSON(ret)
}

func heatmapHandler(c *fiber.Ctx) error {

	resp, err := influx.Read(`SELECT mean("max") AS "mean_max", mean("pcnt") AS "mean_pcnt", mean("people") AS "mean_people" ` +
		`FROM "PureGym"."alltime"."gyms" ` +
		`WHERE "gym" = 'Fareham' ` +
		`GROUP BY time(10m) ` +
		`FILL(0)`,
	)
	if err != nil {
		return err
	}

	var hc = map[string][]int{}

	if len(resp.Results) > 0 && len(resp.Results[0].Series) > 0 {

		var series = resp.Results[0].Series[0]

		for k := range series.Columns {
			if k > 0 {

				for _, vv := range series.Values {

					_, err := time.Parse(time.RFC3339, vv[0].(string))
					if err != nil {
						logger.Error("casting", zap.Error(err))
						continue
					}

					_, err = vv[k].(json.Number).Float64()
					if err != nil {
						logger.Error("casting", zap.Error(err))
						continue
					}

					hc["x"] = append(hc["x"], 1)
				}
			}
		}
	}

	return c.JSON(hc)
}

func submitHandler(c *fiber.Ctx) error {
	return c.SendString("new gym")
}
