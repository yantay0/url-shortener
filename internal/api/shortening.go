package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/yantay0/url-shortener/internal/model"
	"github.com/yantay0/url-shortener/internal/storage"
	"github.com/yantay0/url-shortener/internal/validator"
)

func (app *App) CreateShorterningHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OriginalURL string `json:"original_shorterning"`
		Identifier  string `json:"identifier"`
		UserID      int64  `json:"user_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	shorterning := &model.Shortening{
		OriginalURL: input.OriginalURL,
		Identifier:  input.Identifier,
		UserID:      input.UserID,
	}

	err = app.Storage.Shortenings.Insert(shorterning)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("api/v1/shorternings/%s", shorterning.Identifier))
	fmt.Fprintf(w, "%+v\n", input)
}

func (app *App) ListShorterningsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OriginalURL string
		model.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.OriginalURL = app.readString(qs, "original_url", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "identifier")
	input.Filters.SortSafelist = []string{"original_url", "identifier", "-original_url", "-identifier"}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	shorternings, metadata, err := app.Storage.Shortenings.GetAll(input.OriginalURL, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"shorternings": shorternings, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) ShowShorterningHandler(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	identifier := app.readString(qs, "identifier", "")

	shorterning, err := app.Storage.Shortenings.Get(identifier)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"shorterning": shorterning}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) UpdateShorterningHandler(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	identifier := app.readString(qs, "identifier", "")
	shorterning, err := app.Storage.Shortenings.Get(identifier)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrEditConflict):
			app.editConflictResponse(w, r)
		case errors.Is(err, storage.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	var input struct {
		OriginalURL *string `json:"original_url"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.OriginalURL != nil {
		shorterning.OriginalURL = *input.OriginalURL
	}

	err = app.Storage.Shortenings.Update(shorterning)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"shorterning": shorterning}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) DeleteShorterningHandler(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	Identifier := app.readString(qs, "identifier", "")

	err := app.Storage.Shortenings.Delete(Identifier)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "success.shorterning is deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
