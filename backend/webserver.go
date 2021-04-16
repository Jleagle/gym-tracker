package main

import (
	"encoding/json"
	"time"

	"github.com/Jleagle/puregym-tracker/config"
	"github.com/Jleagle/puregym-tracker/influx"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"go.uber.org/zap"
)

func webserver() error {

	app := fiber.New()

	// Middleware
	app.Use(cache.New(cache.Config{Expiration: time.Minute, KeyGenerator: func(c *fiber.Ctx) string { return c.OriginalURL() }}))
	app.Use(compress.New(compress.Config{Level: compress.LevelBestSpeed}))

	// Routes
	app.Get("/", rootHandler)
	app.Get("/people.json", peopleHandler)
	app.Post("/submit", submitHandler)

	// Serve
	return app.Listen("0.0.0.0:" + config.PortBackend)
}

func rootHandler(c *fiber.Ctx) error {
	return c.SendString("OK")
}

func peopleHandler(c *fiber.Ctx) error {

	var groupBy = c.Query("group")

	switch groupBy {
	case "yearDay", "monthDay", "weekDay", "weekHour", "hour":
	default:
		return nil
	}

	var ret = map[string]map[string][]json.Number{}

	resp, err := influx.Read(`SELECT mean("people") AS "mean_people", mean("pcnt") AS "mean_pcnt" FROM "PureGym"."alltime"."gyms" WHERE time > now()-365d GROUP BY ` + groupBy + ` FILL(0)`)
	if err != nil {
		logger.Error("failed to query influx", zap.Error(err))
		return c.JSON(ret)
	}

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

func submitHandler(c *fiber.Ctx) error {
	return c.SendString("new gym")
}
