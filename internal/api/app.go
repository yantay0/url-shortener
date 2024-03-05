package api

import (
	"log/slog"

	"github.com/yantay0/url-shortener/internal/config"
	"github.com/yantay0/url-shortener/internal/storage"
)

type App struct {
	Config  config.Config
	Logger  *slog.Logger
	Storage storage.Storage
}

func NewApp(cfg config.Config, logger *slog.Logger, storage storage.Storage) *App {
	return &App{
		Config:  cfg,
		Logger:  logger,
		Storage: storage,
	}
}
