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
		// Identifier  string `json:"identifier"`
		// UserID      int64  `json:"user_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	shorten := model.GenerateShortening()

	baseURL := fmt.Sprintf("http://%s:%s", app.Config.IpAdress, app.Config.HTTPServer.Port)
	// http: //localhost:8085/api/v1

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
	fmt.Fprintf(w, "%+v\n", shortURL)
}

// func (app *App) CreateShorterningHandler(w http.ResponseWriter, r *http.Request) {
// 	var input struct {
// 		OriginalURL string `json:"original_shorterning"`
// 		Identifier  string `json:"identifier"`
// 		UserID      int64  `json:"user_id"`
// 	}

// 	err := app.readJSON(w, r, &input)
// 	if err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}

// 	shorterning := &model.Shortening{
// 		OriginalURL: input.OriginalURL,
// 		Identifier:  input.Identifier,
// 		UserID:      input.UserID,
// 	}

// 	err = app.Storage.Shortenings.Insert(shorterning)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}
// 	headers := make(http.Header)
// 	headers.Set("Location", fmt.Sprintf("api/v1/shorternings/%s", shorterning.Identifier))
// 	fmt.Fprintf(w, "%+v\n", input)
// }
