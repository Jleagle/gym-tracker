package main

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/gobuffalo/packr/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"go.uber.org/zap"
)

func webserver() {

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		// AllowOrigins: "https://pgt.jimeagle.com/",
		AllowMethods: "GET",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Use(cache.New(cache.Config{
		Expiration:   10 * time.Minute,
		CacheControl: true,
	}))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))
	app.Use(filesystem.New(filesystem.Config{
		Root:   http.Dir("./public"),
		MaxAge: 3600,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	err := app.Listen("0.0.0.0:" + os.Getenv("PURE_PORT"))
	if err != nil {
		logger.Error("serving webserver", zap.Error(err))
	}
}

type homeTemplate struct {
	Gym string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")

	data := homeTemplate{
		Gym: chi.URLParam(r, "gym"),
	}

	err := template.
		Must(template.ParseFiles("./templates/home.gohtml")).
		ExecuteTemplate(w, "home", data)

	if err != nil {
		logger.Error("executing template", zap.Error(err))
	}
}

type errorTemplate struct {
	Code int
}

func errorHandler(w http.ResponseWriter, _ *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(404)

	data := errorTemplate{
		Code: 404,
	}

	err := template.
		Must(template.ParseFiles("./templates/error.gohtml")).
		ExecuteTemplate(w, "home", data)

	if err != nil {
		logger.Error("executing template", zap.Error(err))
	}
}

func newGymHandler(w http.ResponseWriter, _ *http.Request) {

	_, _ = w.Write([]byte("new gym"))
}

func assetHandler(box *packr.Box, path string) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		b, err := box.Find(strings.TrimPrefix(r.URL.Path, path))
		if err != nil {
			errorHandler(w, r)
			return
		}

		types := map[string]string{
			".js":  "text/javascript",
			".css": "text/css",
			".png": "image/png",
			".jpg": "image/jpeg",
		}

		if val, ok := types[filepath.Ext(r.URL.Path)]; ok {
			w.Header().Add("Content-Type", val)
		}

		_, err = w.Write(b)
		if err != nil {
			logger.Error("writing asset to response", zap.Error(err))
		}
	}
}
