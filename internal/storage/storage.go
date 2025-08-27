package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
	"os"
)

type PGStorage struct {
	Connection *pgx.Conn
}

func New(ctx context.Context, databaseURL string) Storage {
	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return &PGStorage{Connection: conn}
}
