package main

import (
	"context"
	"fmt"
	"github.com/Ocean1342/ya-gofermart/internal/server"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	/* context logger example

	   contextLogger := log.WithFields(log.Fields{
	   		"common": "this is a common field",
	   		"other":  "I also should be logged always",
	   	})

	   	contextLogger.Info("I'll be logged with common and other field")*/
}

func main() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
	defer cancel()
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var name string
	var weight int64
	err = conn.QueryRow(context.Background(), "select name, weight from widgets where id=$1", 42).Scan(&name, &weight)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(name, weight)

	server.Init(ctx)
}
