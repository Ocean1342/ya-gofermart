package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

type Config struct {
	DatabaseURL       string `env:"DATABASE_URI"`
	RunAddr           string `env:"RUN_ADDRESS"`
	AccrualSystemAddr string `env:"ACCULAR_SYSTEM_ADDRESS"`
	SecretKey         string `env:"SECRET_KEY"`
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		DatabaseURL:       os.Getenv("DATABASE_URI"),
		RunAddr:           os.Getenv("RUN_ADDRESS"),
		AccrualSystemAddr: os.Getenv("ACCULAR_SYSTEM_ADDRESS"),
		SecretKey:         os.Getenv("SECRET_KEY"),
	}
}
