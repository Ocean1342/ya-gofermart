package storage

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	"os"
)

type PGStorage struct {
	Connection *pgxpool.Pool
}

func New(ctx context.Context, databaseURL string) Storage {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("could not parse db string: %v", err)
	}

	config.MaxConns = 20
	config.MinConns = 5
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal(err)
	}
	return &PGStorage{Connection: pool}
}
