package main

import (
	"encoding/json"
	"errors"
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
	"github.com/syndtr/goleveldb/leveldb"
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

type Ret struct {
	Group string   `json:"group"`
	Cols  []RetCol `json:"cols"`
}

type RetCol struct {
	X string
	Y map[string]json.Number
}

func newGymHandler(c *fiber.Ctx) error {

	var success bool
	var err error

	defer func() {
		err = c.JSON(map[string]interface{}{"success": success, "message": err.Error()})
		if err != nil {
			logger.Error("returning response", zap.Error(err))
		}
	}()

	var db *leveldb.DB
	db, err = leveldb.OpenFile("./leveldb/", nil)
	if err != nil {
		logger.Error("opening database", zap.Error(err))
		return err
	}

	//goland:noinspection GoUnhandledErrorResult
	defer db.Close()

	// Get data form request
	var request []string
	err = json.Unmarshal(c.Body(), &request)
	if err != nil {
		logger.Error("opening database", zap.Error(err))
		return err
	}

	if len(request) != 2 {
		err = errors.New("invalid post data")
		logger.Error("invalid post data", zap.Error(err))
		return err
	}

	var gym, errorString string
	_, gym, err, errorString = loginAndCheckMembers(request[0], request[1])
	if err != nil {
		logger.Error("scraping", zap.Error(err))
		return err
	}

	if errorString != "" {
		err = errors.New(errorString)
		return nil
	}

	var b []byte
	b, err = db.Get([]byte(gym), nil)
	if err != nil {
		logger.Error("reading gym from db", zap.Error(err))
		return err
	}

	// Update database
	var gyms Credentials
	err = json.Unmarshal(b, &gyms)
	if err != nil {
		logger.Error("unmarshaling gyms", zap.Error(err))
		return err
	}

	gyms[request[0]] = request[1]

	// Save back to file
	b, err = json.Marshal(gyms)
	if err != nil {
		logger.Error("marshaling gyms", zap.Error(err))
		return err
	}

	err = db.Put([]byte(gym), b, nil)
	if err != nil {
		logger.Error("reading gym from db", zap.Error(err))
		return err
	}

	// Return
	success = true
	return err
}

type Credentials map[string]string
