package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *App) Routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/healthcheck", app.HealthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/urls", app.ListUrlsHandler)
	router.HandlerFunc(http.MethodGet, "/urls/:id", app.ShowUrlHandler)
	router.HandlerFunc(http.MethodPost, "/urls", app.CreateUrlHandler)
	router.HandlerFunc(http.MethodPut, "/urls/:id", app.UpdateUrlHandler)
	router.HandlerFunc(http.MethodDelete, "/urls/:id", app.DeleteUrlHandler)
	return router
}
