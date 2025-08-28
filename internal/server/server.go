package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"gofermart/config"
	"gofermart/internal/api"
	"net/http"
)

func Init(cfg *config.Config, handler *api.Handler) {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(loggable)
	r.Post("/api/user/register", handler.UserRegister)
	r.Post("/api/user/login", handler.UserAuth)
	r.With(handler.Authenticate).Post("/api/user/orders", handler.LoadOrderNumber)
	r.With(handler.Authenticate).Get("/api/user/orders", handler.GetOrders)
	r.With(handler.Authenticate).Get("/api/user/balance", handler.GetUserBalance)
	r.With(handler.Authenticate).Post("/api/user/balance/withdraw", handler.UserWithdraw)

	server := &http.Server{
		Addr:    cfg.RunAddr,
		Handler: r,
	}
	err := server.ListenAndServe()
	if err != nil {
		logrus.Errorf("could not start server. err: %s", err)
	}
}
