package main

import (
	"os"

	api "github.com/yantay0/url-shortener/internal/api"
	"github.com/yantay0/url-shortener/internal/config"
	"github.com/yantay0/url-shortener/internal/mailer"
	"github.com/yantay0/url-shortener/internal/storage"

	"github.com/yantay0/url-shortener/internal/lib/logger/jsonlog"
	"github.com/yantay0/url-shortener/internal/storage/postgres"
)

// @title URL-shortener
// @host localhost: 8085

// const (
// 	envLocal = "local"
// 	envDev   = "dev"
// 	envProd  = "prod"
// 	version  = "1.0.0"
// )

func main() {
	cfg := config.MustLoad()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := postgres.OpenDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()
	logger.PrintInfo("database conntection pool established", nil)

	app := api.NewApp(*cfg, logger, storage.New(db), mailer.New(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Sender))

	err = app.Serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
