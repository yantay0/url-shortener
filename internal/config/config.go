package config

import (
	"log"
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string `env:"ENV" env-default:"develepment"`
	StoragePath string `env:"STORAGE_PATH" env-required:"true"`
	HTTPServer  HTTPServer
}

type HTTPServer struct {
	Host        string        `env:"HOST" env-default:"0.0.0.0"`
	Port        string        `env:"PORT" env-default:"8080"`
	Timeout     time.Duration `env:"SERVER_TIMEOUT" env-default:"5s"`
	IdleTimeout time.Duration `env:"IDLE_TIMEOUT" env-default:"60s"`
}

func MustLoad() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("unable to load .env file: %e", err)
	}

	var cfg Config

	err = env.Parse(&cfg)
	if err != nil {
		log.Fatalf("unable to parse ennvironment variables: %e", err)
	}
	return cfg
}
