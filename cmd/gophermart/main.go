package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"gofermart/config"
	"gofermart/internal/api"
	"gofermart/internal/auth"
	"gofermart/internal/server"
	"gofermart/internal/storage"
	"os"
	"time"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	cfg := config.New()
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()
	repo := storage.New(ctx, cfg.DatabaseURL)
	jwtAuth := auth.New(cfg.SecretKey, cfg.TokenTTL)
	handler := api.New(repo, jwtAuth)
	server.Init(cfg, handler)
}
