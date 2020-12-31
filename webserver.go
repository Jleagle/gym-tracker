package main

import (
	"compress/flate"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gobuffalo/packr/v2"
	"go.uber.org/zap"
)

var (
	publicBox = packr.New("public", "./public")
	//templatesBox = packr.New("templates", "./templates")
)

func webserver() {

	r := chi.NewRouter()
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.NewCompressor(flate.DefaultCompression, "text/html", "text/css", "text/javascript", "application/json", "application/javascript").Handler)

	r.Get("/public/*", assetHandler(publicBox, "/public"))

	r.Get("/", homeHandler)
	r.Get("/{gym}", homeHandler)
	r.Post("/", newGymHandler)

	r.NotFound(errorHandler)

	port := os.Getenv("PURE_PORT")
	if port == "" {
		port = "9030"
	}

	s := &http.Server{
		Addr:              "0.0.0.0:" + port,
		Handler:           r,
		ReadTimeout:       2 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	logger.Info("Starting Frontend on " + "http://localhost:" + port)

	err := s.ListenAndServe()
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
