package main

import (
	"compress/flate"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

func webserver() {

	r := chi.NewRouter()
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.NewCompressor(flate.DefaultCompression, "text/html", "text/css", "text/javascript", "application/json", "application/javascript").Handler)

	r.Get("/", pages.HomeHandler)

	// 404
	r.NotFound(pages.Error404Handler)

	s := &http.Server{
		Addr:              "0.0.0.0:" + os.Getenv("PURE_PORT"),
		Handler:           r,
		ReadTimeout:       2 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	logger.Info("Starting Frontend on " + "http://" + s.Addr)

	err := s.ListenAndServe()
	if err != nil {
		logger.Error("serving webserver", zap.Error(err))
	}
}
