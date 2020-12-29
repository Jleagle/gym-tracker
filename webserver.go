package main

import (
	"compress/flate"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gobuffalo/packr/v2"
	"go.uber.org/zap"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	assetBox = packr.New("assets", "./assets")
)

func webserver() {

	r := chi.NewRouter()
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.NewCompressor(flate.DefaultCompression, "text/html", "text/css", "text/javascript", "application/json", "application/javascript").Handler)

	r.Get("/assets/*", assetHandler(assetBox, "/assets"))

	r.Get("/", homeHandler)
	r.Get("/{gym}", homeHandler)
	r.Post("/", newGymHandler)

	r.NotFound(errorHandler)

	s := &http.Server{
		Addr:              "0.0.0.0:" + os.Getenv("PURE_PORT"),
		Handler:           r,
		ReadTimeout:       2 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	logger.Info("Starting Frontend on " + "http://localhost:" + os.Getenv("PURE_PORT"))

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
		Must(template.ParseFiles("./assets/home.gohtml")).
		ExecuteTemplate(w, "home", data)

	if err != nil {
		logger.Error("executing template", zap.Error(err))
	}
}

type errorTemplate struct {
	Code int
}

func errorHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(404)

	data := errorTemplate{
		Code: 404,
	}

	err := template.
		Must(template.ParseFiles("./assets/error.gohtml")).
		ExecuteTemplate(w, "home", data)

	if err != nil {
		logger.Error("executing template", zap.Error(err))
	}
}

func newGymHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("new gym"))
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
