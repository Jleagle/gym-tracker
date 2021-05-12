package main

import (
	"encoding/json"
	"errors"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Jleagle/gym-tracker/config"
	"github.com/Jleagle/gym-tracker/datastore"
	"github.com/Jleagle/gym-tracker/helpers"
	"github.com/Jleagle/gym-tracker/influx"
	"github.com/Jleagle/gym-tracker/log"
	"github.com/gofiber/fiber/v2"
	fiverCache "github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

func webserver() error {

	app := fiber.New()

	// Middleware
	if config.Environment == "PRODUCTION" {
		app.Use(fiverCache.New(fiverCache.Config{Expiration: time.Minute, KeyGenerator: func(c *fiber.Ctx) string { return c.OriginalURL() }}))
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

var gc = cache.New(time.Minute, time.Second*10)

func peopleHandler(c *fiber.Ctx) error {

	var groupBy = c.Query("group")
	var ret = Ret{Group: groupBy, Cols: []RetCol{}}
	var query string
	var key = "influx-" + groupBy

	if config.Environment == "PRODUCTION" {
		cached, found := gc.Get(key)
		if found {
			return c.JSON(cached)
		}
	}

	switch groupBy {
	case "yearDay", "monthDay", "weekDay", "weekHour", "dayHour":

		query = `SELECT mean("people") AS "members", mean("pcnt") AS "percent" FROM "GymTracker"."alltime"."gyms" WHERE time > now()-365d GROUP BY ` + groupBy + ` FILL(0)`

	case "now":

		query = `SELECT mean("people") AS "members", mean("pcnt") AS "percent" FROM "GymTracker"."alltime"."gyms" WHERE time > now()-24h GROUP BY time(10m) FILL(0)`

	default:
		return c.JSON(ret)
	}

	resp, err := influx.Read(query)
	if err != nil {
		log.Instance.Error("failed to query influx", zap.Error(err))
		return c.JSON(ret)
	}

	for _, result := range resp.Results {
		for _, series := range result.Series {
			for _, row := range series.Values {

				var x string

				if groupBy == "now" {

					t, err := time.Parse(time.RFC3339, row[0].(string))
					if err != nil {
						log.Instance.Error("parsing time", zap.Error(err))
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

	gc.Set(key, ret, cache.DefaultExpiration)

	return c.JSON(ret)
}

type Ret struct {
	Group string   `json:"group"`
	Cols  []RetCol `json:"cols"`
}

type RetCol struct {
	X string
	Y map[string]json.Number
}

//goland:noinspection GoErrorStringFormat
func newGymHandler(c *fiber.Ctx) error {

	var success bool
	var err error

	defer func() {

		resp := map[string]interface{}{"success": success}
		if err != nil {
			resp["message"] = err.Error()
		}

		err = c.JSON(resp)
		if err != nil {
			log.Instance.Error("returning response", zap.Error(err))
		}
	}()

	// Get data form request
	var request []string
	err = json.Unmarshal(c.Body(), &request)
	if err != nil {
		log.Instance.Error("opening database", zap.Error(err))
		return err
	}

	if len(request) != 2 {
		err = errors.New("invalid post data")
		log.Instance.Error("invalid post data", zap.Error(err))
		return err
	}

	request[0] = strings.TrimSpace(request[0])
	request[1] = strings.TrimSpace(request[1])

	if //goland:noinspection RegExpRedundantEscape
	!regexp.MustCompile(`.+\@.+\..+`).Match([]byte(request[0])) {
		err = errors.New("Invalid Email")
		return nil
	}

	if !regexp.MustCompile("^[0-9]{8}$").Match([]byte(request[1])) {
		err = errors.New("Invalid PIN")
		return nil
	}

	_, gym, err, errorString := scrape(datastore.Credential{Email: request[0], PIN: request[1]})
	if err != nil {
		log.Instance.Error("scraping", zap.Error(err))
		return err
	}

	if errorString != "" {
		err = errors.New(errorString)
		return nil
	}

	err = datastore.SaveNewCredential(request[0], request[1], gym)
	if err != nil {
		log.Instance.Error("saving to datastore", zap.Error(err))
		return err
	}

	// Return
	success = true
	return err
}
