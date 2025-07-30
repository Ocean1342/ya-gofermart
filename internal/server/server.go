package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Init(ctx context.Context) {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(loggable)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	err := server.ListenAndServe()
	if err != nil {
		logrus.Errorf("could not start server. err: %s", err)
	}
}
