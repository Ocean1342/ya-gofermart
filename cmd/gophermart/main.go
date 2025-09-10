package main

import (
	"context"
	_ "github.com/jackc/pgx/v5/stdlib"
	log "github.com/sirupsen/logrus"
	"gofermart/config"
	accrualsystem "gofermart/internal/accrual-system"
	"gofermart/internal/api"
	"gofermart/internal/auth"
	"gofermart/internal/migrator"
	orderprocessor "gofermart/internal/order-processor"
	"gofermart/internal/server"
	"gofermart/internal/service"
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
	ctx := context.Background()
	repo := storage.New(ctx, cfg.DatabaseURL)
	err := migrator.Migrate(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("could not migrate database: %v", err)
	}
	system := accrualsystem.New(cfg.AccrualSystemAddr, "api/orders")
	orderProcessor := orderprocessor.New(system, repo, 60*time.Second)
	go func() {
		orderProcessor.Process(ctx)
	}()
	handler := api.New(repo, auth.New(cfg.SecretKey, cfg.TokenTTL), orderProcessor, service.New(repo))
	server.Init(cfg, handler)
}
