package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const BASE_URL = "/api/v1"

func (app *App) Routes() *httprouter.Router {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, BASE_URL+"/healthcheck", app.HealthcheckHandler)
	router.HandlerFunc(http.MethodGet, BASE_URL+"/urls", app.ListUrlsHandler)
	router.HandlerFunc(http.MethodGet, BASE_URL+"/urls/:id", app.ShowUrlHandler)
	router.HandlerFunc(http.MethodPost, BASE_URL+"/urls", app.CreateUrlHandler)
	router.HandlerFunc(http.MethodPatch, BASE_URL+"/urls/:id", app.UpdateUrlHandler)
	router.HandlerFunc(http.MethodDelete, BASE_URL+"/urls/:id", app.DeleteUrlHandler)
	return router
}
