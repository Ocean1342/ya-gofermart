package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseURL       string        `env:"DATABASE_URI"`
	RunAddr           string        `env:"RUN_ADDRESS"`
	AccrualSystemAddr string        `env:"ACCULAR_SYSTEM_ADDRESS"`
	SecretKey         string        `env:"SECRET_KEY"`
	TokenTTL          time.Duration `env:"TOKEN_TTL"`
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Errorf("Error loading .env file")
	}

	ttl, err := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	if err != nil {
		log.Errorf("could not set token ttl.err:%s", err)
	}
	if ttl == 0 {
		ttl = 24
	}
	addr := os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	log.Infof("ACCRUAL_SYSTEM_ADDRESS: %s", addr)
	return &Config{
		DatabaseURL:       os.Getenv("DATABASE_URI"),
		RunAddr:           os.Getenv("RUN_ADDRESS"),
		AccrualSystemAddr: addr,
		SecretKey:         os.Getenv("SECRET_KEY"),
		TokenTTL:          time.Duration(ttl) * time.Hour,
	}
}
