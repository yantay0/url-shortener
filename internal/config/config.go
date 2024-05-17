package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	DB         `yaml:"db"`
	HTTPServer `yaml:"http_server"`
	SMTP       `yaml:"smtp"`
	Limiter    `yaml:"limiter"`
}

type SMTP struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Sender   string `yaml:"sender"`
}

type DB struct {
	Dsn          string `yaml:"dsn" env-required:"true"`
	MaxOpenConns int    `yaml:"maxOpenConns" env-default:"25"`
	MaxIdleConns int    `yaml:"maxIdleConns" env-default:"25"`
	MaxIdleTime  string `yaml:"maxIdleTime" env-default:"15m"`
}

type HTTPServer struct {
	IpAdress    string        `yaml:"ip_address" env-default:"localhost"`
	Port        string        `yaml:"port" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Limiter struct {
	RPS     float64 `yaml:"rps" env-default:"2"` // Rate limiter maximum requests per second
	Burst   int     `yaml:"burst" env-default:"4"`
	Enabled bool    `yaml:"enabled" env-default:"true"`
}

func MustLoad() *Config {
	configPath := "./internal/config/prod.yaml"
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
