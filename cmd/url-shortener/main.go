package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	api "github.com/yantay0/url-shortener/internal/api"
	"github.com/yantay0/url-shortener/internal/config"
	"github.com/yantay0/url-shortener/internal/storage"

	"github.com/yantay0/url-shortener/internal/lib/logger/handler/slogpretty"
	"github.com/yantay0/url-shortener/internal/lib/logger/sl"
	"github.com/yantay0/url-shortener/internal/storage/postgres"
)

// @title URL-shortener
// @host localhost: 8085

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	log = log.With(slog.String("env", cfg.Env)) // current env is added to each log
	log.Debug("debug messages are enabled")

	db, err := postgres.OpenDB(cfg)
	if err != nil {
		log.Error("failed to open db connection", sl.Err(err))
		os.Exit(1)
	}

	app := api.NewApp(*cfg, log, storage.NewModels(db))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HTTPServer.Port),
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	defer db.Close()

	log.Info("starting server", cfg.Env, cfg.HTTPServer.Port)
	err = srv.ListenAndServe()
	log.Error("failed to start server", sl.Err(err))
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch envLocal {
	case envLocal:
		log = setupPrettySlog()

	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(

			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
