package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/yantay0/url-shortener/internal/config"
	"github.com/yantay0/url-shortener/internal/data"
	sl "github.com/yantay0/url-shortener/internal/lib/logger"
	"github.com/yantay0/url-shortener/internal/repository/postgres"
)

// @title URL-shortener
// @host localhost: 8085

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type application struct {
	config config.Config
	logger *slog.Logger
	models data.Models
}

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	log = log.With(slog.String("env", cfg.Env)) // current env is added to each log
	log.Info("starting app")
	log.Debug("debug messages are enabled")

	repository, err := postgres.OpenDB(cfg)
	if err != nil {
		log.Error("failed to open db connection", sl.Err(err))
		os.Exit(1)
	}

	defer repository.Db.Close()

	app := &application{
		config: *cfg,
		logger: log,
		models: data.NewModels(repository.Db),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      mux,
		IdleTimeout:  cfg.IdleTimeout,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Info("starting %s server on %d", cfg.Env, cfg.HTTPServer.Port)
	// err := srv.ListenAndServe()
	err = srv.ListenAndServe()
	log.Error("failed to start server", sl.Err(err))

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
