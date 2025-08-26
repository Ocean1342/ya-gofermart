package main

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
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
	migrate(cfg.DatabaseURL)
	jwtAuth := auth.New(cfg.SecretKey, cfg.TokenTTL)
	handler := api.New(repo, jwtAuth)
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
