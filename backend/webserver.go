package main

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/Jleagle/puregym-tracker/config"
	"github.com/Jleagle/puregym-tracker/influx"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"go.uber.org/zap"
)

func webserver() error {

	app := fiber.New()

	// Middleware
	// app.Use(cache.New(cache.Config{Expiration: time.Minute, KeyGenerator: func(c *fiber.Ctx) string { return c.OriginalURL() }}))
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
	var ret = Ret{Group: groupBy, Cols: []RetCol{}}
	var query string

	switch groupBy {
	case "yearDay", "monthDay", "weekDay", "weekHour", "hour":

		query = `SELECT mean("people") AS "members", mean("pcnt") AS "percent" FROM "PureGym"."alltime"."gyms" WHERE time > now()-365d GROUP BY ` + groupBy + ` FILL(0)`

	case "now":

		query = `SELECT mean("people") AS "members", mean("pcnt") AS "percent" FROM "PureGym"."alltime"."gyms" WHERE time > now()-24h GROUP BY time(10m) FILL(0)`

	default:
		return c.JSON(ret)
	}

	resp, err := influx.Read(query)
	if err != nil {
		logger.Error("failed to query influx", zap.Error(err))
		return c.JSON(ret)
	}

	for _, result := range resp.Results {
		for _, series := range result.Series {
			for _, row := range series.Values {

				var x string

				if groupBy == "now" {

					t, err := time.Parse(time.RFC3339, row[0].(string))
					if err != nil {
						logger.Error("parsing time", zap.Error(err))
					}

					x = strconv.FormatInt(t.Unix(), 10)

				} else {
					x = series.Tags[groupBy]
				}

				y := map[string]json.Number{}
				for k, col := range row {
					if k > 0 {
						y[series.Columns[k]] = col.(json.Number)
					}
				}

				ret.Cols = append(ret.Cols, RetCol{X: x, Y: y})
			}
		}
	}

	return c.JSON(ret)
}

func submitHandler(c *fiber.Ctx) error {
	return c.SendString("new gym")
}

type Ret struct {
	Group string   `json:"group"`
	Cols  []RetCol `json:"cols"`
}

type RetCol struct {
	X string
	Y map[string]json.Number
}
