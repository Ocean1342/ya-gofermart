package main

import (
	"context"
	"github.com/Ocean1342/ya-gofermart/config"
	"github.com/Ocean1342/ya-gofermart/internal/server"
	"github.com/Ocean1342/ya-gofermart/internal/storage"
	log "github.com/sirupsen/logrus"
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

	server.Init(ctx, cfg, repo)
}
