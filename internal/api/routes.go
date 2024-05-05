package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const BASE_URL = "/api/v1"

func (app *App) Routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, BASE_URL+"/healthcheck", app.HealthcheckHandler)

	router.HandlerFunc(http.MethodGet, BASE_URL+"/urls", app.requirePermission("urls:read", app.ListUrlsHandler))
	router.HandlerFunc(http.MethodGet, BASE_URL+"/urls/:id", app.requirePermission("urls:read", app.ShowUrlHandler))
	router.HandlerFunc(http.MethodPost, BASE_URL+"/urls", app.requirePermission("urls:write", app.CreateUrlHandler))
	router.HandlerFunc(http.MethodPatch, BASE_URL+"/urls/:id", app.requirePermission("urls:write", app.UpdateUrlHandler))
	router.HandlerFunc(http.MethodDelete, BASE_URL+"/urls/:id", app.requirePermission("urls:write", app.DeleteUrlHandler))

	router.HandlerFunc(http.MethodPost, BASE_URL+"/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, BASE_URL+"/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, BASE_URL+"/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
