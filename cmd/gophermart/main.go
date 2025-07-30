package main

import (
	"context"
	"github.com/Ocean1342/ya-gofermart/internal/server"
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
	//TODO: graceful shutdown
	server.Init(ctx)
}
