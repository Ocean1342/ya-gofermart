package server

import (
	"context"
	"github.com/Ocean1342/ya-gofermart/config"
	"github.com/Ocean1342/ya-gofermart/internal/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Init(ctx context.Context, cfg *config.Config, handler *api.Handler) {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(loggable)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Post("/api/user/register", handler.UserRegister)
	r.Post("/api/user/login", handler.UserAuth)
	r.With(authable).Post("/api/user/orders", handler.LoadOrderNumber)
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	err := server.ListenAndServe()
	if err != nil {
		logrus.Errorf("could not start server. err: %s", err)
	}
}
