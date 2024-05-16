package api

import (
	"net/http"
)

func (app *App) ListUserUrlsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	urls, err := app.Storage.Shortenings.GetUserAllShortenings(userID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"shortenings": urls}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
