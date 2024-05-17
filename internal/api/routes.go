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

	router.HandlerFunc(http.MethodGet, BASE_URL+"/shortenings", app.requirePermission("shortenings:read", app.ListShorterningsHandler))
	router.HandlerFunc(http.MethodGet, BASE_URL+"/shortenings/:id", app.requirePermission("shortenings:read", app.ShowShorterningHandler))
	router.HandlerFunc(http.MethodPost, BASE_URL+"/shortenings", app.requirePermission("shortenings:write", app.CreateShorterningHandler))
	router.HandlerFunc(http.MethodPatch, BASE_URL+"/shortenings/:id", app.requirePermission("shortenings:write", app.UpdateShorterningHandler))
	router.HandlerFunc(http.MethodDelete, BASE_URL+"/shortenings/:id", app.requirePermission("shortenings:write", app.DeleteShorterningHandler))

	router.HandlerFunc(http.MethodPost, BASE_URL+"/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, BASE_URL+"/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodGet, BASE_URL+"/users/:id/shortenings", app.requirePermission("shortenings:read", app.listUserShorteningsHandler))
	router.HandlerFunc(http.MethodPost, BASE_URL+"/users/:id/shortenings", app.createShorteningFromURLHandler)

	router.HandlerFunc(http.MethodPost, BASE_URL+"/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
