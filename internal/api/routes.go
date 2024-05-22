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

	// Define specific routes first, including those starting with /api/v1/
	router.HandlerFunc(http.MethodGet, BASE_URL+"/healthcheck", app.HealthcheckHandler)

	router.HandlerFunc(http.MethodGet, BASE_URL+"/shortenings", app.requirePermission("shortenings:read", app.ListShorterningsHandler))
	router.HandlerFunc(http.MethodGet, BASE_URL+"/shortenings/:identifier", app.requirePermission("shortenings:read", app.ShowShorterningHandler))
	router.HandlerFunc(http.MethodPatch, BASE_URL+"/shortenings/:identifier", app.requirePermission("shortenings:write", app.UpdateShorterningHandler))
	router.HandlerFunc(http.MethodDelete, BASE_URL+"/shortenings/:identifier", app.requirePermission("shortenings:write", app.DeleteShorterningHandler))

	router.HandlerFunc(http.MethodPost, BASE_URL+"/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, BASE_URL+"/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodGet, BASE_URL+"/users/:id/shortenings", app.requirePermission("shortenings:read", app.listUserShorteningsHandler))
	router.HandlerFunc(http.MethodPost, BASE_URL+"/users/:id/shortenings", app.createShorteningFromURLHandler)
	router.HandlerFunc(http.MethodPost, BASE_URL+"/tokens/authentication", app.createAuthenticationTokenHandler)

	// Define the wildcard route after all specific routes
	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}

// func (app *App) Routes() http.Handler {
// 	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
// 	router.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

// 	router.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

// 	// healthcheck
// 	router.HandleFunc("/healthcheck", app.HealthcheckHandler).Methods("GET")

// 	router.HandleFunc("/shortenings", app.requirePermission("shortenings:read", app.ListShorterningsHandler)).Methods("GET")
// 	router.HandleFunc("/shortenings", app.requirePermission("shortenings:write", app.CreateShorterningHandler)).Methods("POST")
// 	router.HandleFunc("/shortenings/identifier}", app.ShowShorterningHandler).Methods("GET")
// 	router.HandleFunc("/shortenings/identifier}", app.requirePermission("shortenings:write", app.UpdateShorterningHandler)).Methods("PUT")
// 	router.HandleFunc("/shortenings/identifier}", app.requirePermission("shortenings:write", app.DeleteShorterningHandler)).Methods("DELETE")

// 	// User handlers with Authentication
// 	router.HandleFunc("/users", app.requirePermission("shortenings:read", app.registerUserHandler)).Methods("POST")
// 	router.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")
// 	router.HandleFunc("/tokens/authentication", app.createAuthenticationTokenHandler).Methods("POST")

// 	router.HandleFunc("/users/:id/shortenings", app.requirePermission("shortenings:read", app.listUserShorteningsHandler)).Methods("GET")
// 	router.HandleFunc("/users/:id/shortenings", app.createShorteningFromURLHandler).Methods("POST")

// 	router.HandleFunc("/{identifier}", app.requirePermission("shortenings:read", app.redirectHandler)).Methods("GET")

// 	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
// }
