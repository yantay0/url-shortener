package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yantay0/url-shortener/internal/model"
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
	var input struct {
		OriginalURL string `json:"original_url"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	shorten := model.GenerateShortening()

	baseURL := fmt.Sprintf("http://%s:%s", app.Config.IpAdress, app.Config.HTTPServer.Port)

	shorterning := &model.Shortening{
		OriginalURL: input.OriginalURL,
	}

	shortURL, err := model.PrependBaseURL(baseURL, shorten)
	if err != nil {
		log.Printf("error generating full URL: %v", err)
		app.badRequestResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("api/v1/shorternings/%s", shorterning.Identifier))

	err = app.writeJSON(w, http.StatusCreated, envelope{"short_url": shortURL}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
