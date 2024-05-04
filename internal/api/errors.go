package api

import (
	"fmt"
	"net/http"
)

func (app *App) logError(r *http.Request, err error) {
	app.Logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (app *App) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
	}
}

func (app *App) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *App) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// methodNotAllowedResponse method is used to send a 405 Method Not Allowed status code and
// JSON response to the client.
func (app *App) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// badRequestResponse sends JSON-formatted error message with 400 Bad Request status code.
func (app *App) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// failedValidationResponse sends JSON-formatted error message to client with UnprocessableEntity
// 422 status code when Validation fails.
// Note that the errors parameter here has the type map[string]string,
// which is exact the same as the errors map contained in our Validator type.
func (app *App) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

// editConflictResponse sends a JSON-formatted error message to the client with a 409 Conflict
// status code.
func (app *App) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, message)
}

// rateLimitExceedResponse sends a JSON-formatted error message with a 429 Too Many Requests
// status code to the client.
func (app *App) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limited exceeded"
	app.errorResponse(w, r, http.StatusTooManyRequests, message)
}

// invalidCredentialsResponse sends a JSON-formatted error with a 401 Unauthorized status code
// to the client.
func (app *App) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

// invalidAuthenticationTokenResponse sends a JSON-formatted error with a 401 Unauthorized status
// code and "WWW-Authenticate: Bearer" header to the client.
func (app *App) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

// authenticationRequiredResponse sends a JSON-formatted error with a 401 Unauthorized status code
// to the client.
func (app *App) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

// inactiveAccountResponse sends a JSON-formatted error with a 403 Forbidden status code to the
// client.
func (app *App) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account must be activated to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}

func (app *App) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account doesn't have the necessary permissions to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}
