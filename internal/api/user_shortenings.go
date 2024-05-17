package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yantay0/url-shortener/internal/model"
	"github.com/yantay0/url-shortener/internal/storage"
	"github.com/yantay0/url-shortener/internal/validator"
)

func (app *App) listUserShorteningsHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *App) createShorteningFromURLHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var input struct {
		OriginalURL string `json:"original_url"`
		Identifier  string `json:"identifier,omitempty"` // Make Identifier optional
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Check if Identifier is provided in the request
	if input.Identifier == "" {
		// If not provided, generate a new identifier server-side
		shorten := model.GenerateShortening()
		input.Identifier = shorten
	}

	baseURL := fmt.Sprintf("http://%s:%s", app.Config.HTTPServer.IpAdress, app.Config.HTTPServer.Port) // Fixed typo in IpAddress

	shortURL, err := model.PrependBaseURL(baseURL, input.Identifier)
	if err != nil {
		log.Printf("error generating full URL: %v", err)
		app.badRequestResponse(w, r, err)
		return
	}

	shortening := &model.Shortening{
		Identifier:  input.Identifier,
		OriginalURL: input.OriginalURL,
		UserID:      userID,
	}
	v := validator.New()
	err = app.Storage.Shortenings.SaveUserShortening(shortening)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrIdentifierExists):
			v.AddError("message", "identifier already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("api/v1/shorternings/%s", shortening.Identifier))

	err = app.writeJSON(w, http.StatusCreated, envelope{"short_url": shortURL}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *App) redirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	identifier := vars["identifier"]

	shortening, err := app.Storage.Shortenings.GetOriginalUrl(identifier)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	http.Redirect(w, r, shortening.OriginalURL, http.StatusMovedPermanently)
}
