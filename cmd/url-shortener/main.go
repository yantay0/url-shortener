package main

import (
	"log/slog"
	"os"

	"github.com/yantay0/url-shortener/internal/config"
	"github.com/yantay0/url-shortener/internal/repository/postgres"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env)) // current env is added to each log

	log.Info("starting app")
	log.Debug("debug messages are enabled")

	repository, err := postgres.OpenDB(cfg)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	_ = repository

	// 	db, err := openDB(cfg)
	// 	if err != nil {
	// 		log.Error(err.Error())
	// 	}

	// defer func() {
	// 	if err := db.Close(); err != nil {
	// 		log.Error(err.Error())
	// 	}
	// }()

}
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

		// envDev, envProd
	}
	return log
}
