package api

import (
	"net/http"

	"github.com/yantay0/url-shortener/internal/config"
	"github.com/yantay0/url-shortener/internal/lib/logger/jsonlog"
	"github.com/yantay0/url-shortener/internal/mailer"
	"github.com/yantay0/url-shortener/internal/storage"
)

type App struct {
	Config  config.Config
	Logger  *jsonlog.Logger
	Storage storage.Storage
	Mailer  mailer.Mailer
}

func NewApp(cfg config.Config, logger *jsonlog.Logger, storage storage.Storage, mailer mailer.Mailer) *App {
	return &App{
		Config:  cfg,
		Logger:  logger,
		Storage: storage,
		Mailer:  mailer,
	}
}

func (app *App) catchAllRequestHandler(w http.ResponseWriter, r *http.Request, path string) {
	// Log the request or perform other actions
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not Found"))
}
