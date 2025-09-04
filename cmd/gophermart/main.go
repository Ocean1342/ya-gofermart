package main

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	log "github.com/sirupsen/logrus"
	"gofermart/config"
	accrual_system "gofermart/internal/accrual-system"
	"gofermart/internal/api"
	"gofermart/internal/auth"
	order_processor "gofermart/internal/order-processor"
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
	ctx := context.Background()
	repo := storage.New(ctx, cfg.DatabaseURL)
	migrate(cfg.DatabaseURL)
	jwtAuth := auth.New(cfg.SecretKey, cfg.TokenTTL)
	system := accrual_system.New(cfg.AccrualSystemAddr, "api/orders")
	orderProcessor := order_processor.New(system, repo, 2*time.Second)
	go func() {
		orderProcessor.Process(ctx)
	}()
	handler := api.New(repo, jwtAuth, orderProcessor)
	server.Init(cfg, handler)
}

func migrate(dbURL string) {
	db, err := sql.Open("pgx", dbURL)
	defer func() {
		err = db.Close()
		if err != nil {
			log.Errorf("could not close db connection:%s", err)
		}
	}()
	if err != nil {
		panic("could not run migration")
	}
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
	if err := goose.Up(db, "./migrations"); err != nil {
		panic(err)
	}
}
