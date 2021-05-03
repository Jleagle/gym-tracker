package main

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Jleagle/gym-tracker/config"
	"github.com/Jleagle/gym-tracker/helpers"
	"github.com/Jleagle/gym-tracker/influx"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

func webserver() error {

	app := fiber.New()

	// Middleware
	if config.Environment == "PRODUCTION" {
		app.Use(cache.New(cache.Config{Expiration: time.Minute, KeyGenerator: func(c *fiber.Ctx) string { return c.OriginalURL() }}))
	}
	app.Use(compress.New(compress.Config{Level: compress.LevelBestSpeed}))
	app.Use(cors.New(cors.Config{AllowOrigins: "*", AllowMethods: "GET POST"})) // for the new gym form

	// Routes
	app.Get("/", rootHandler)
	app.Get("/people.json", peopleHandler)
	app.Post("/new-gym", newGymHandler)

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

					// Move Sunday from first to last
					if strings.HasPrefix(groupBy, "week") && strings.HasPrefix(x, "0") {
						x = helpers.ReplaceAtIndex(x, '7', 0)
					}
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

	sort.SliceStable(ret.Cols, func(i, j int) bool {

		pieces1 := strings.Split(ret.Cols[i].X, "-")
		pieces2 := strings.Split(ret.Cols[j].X, "-")

		i1, err1 := strconv.Atoi(pieces1[0])
		i2, err2 := strconv.Atoi(pieces2[0])
		if err1 == nil && err2 == nil {

			if i1 == i2 && len(pieces1) > 1 && len(pieces2) > 1 {

				i1, err1 = strconv.Atoi(pieces1[1])
				i2, err2 = strconv.Atoi(pieces2[1])
				if err1 == nil && err2 == nil {
					return i1 < i2 // Sort by second value
				}
			}

			return i1 < i2 // Sort by first value
		}

		return pieces1[0] < pieces2[0] // Alphabetically
	})

	return c.JSON(ret)
}

func newGymHandler(c *fiber.Ctx) error {
	return c.JSON(map[string]interface{}{
		"success": 1,
		"message": "OK",
	})
}

type Ret struct {
	Group string   `json:"group"`
	Cols  []RetCol `json:"cols"`
}

type RetCol struct {
	X string
	Y map[string]json.Number
}
