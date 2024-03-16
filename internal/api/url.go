package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/yantay0/url-shortener/internal/model"
	"github.com/yantay0/url-shortener/internal/storage"
	"github.com/yantay0/url-shortener/internal/validator"
)

func (app *App) CreateUrlHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OriginalUrl string `json:"original_url"`
		ShortUrl    string `json:"short_url"`
	}

	err := app.readJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	url := &model.Url{
		OriginalUrl: input.OriginalUrl,
		ShortUrl:    input.ShortUrl,
	}

	err = app.Storage.Urls.Insert(url)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("api/v1/urls/%d", url.ID))
	fmt.Fprintf(w, "%+v\n", input)
}

func (app *App) ListUrlsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OriginalUrl string
		ShortUrl    string
		Sort        string
	}

	v := validator.New()
	qs := r.URL.Query()

	input.OriginalUrl = app.readString(qs, "original_url", "")
	input.ShortUrl = app.readString(qs, "short_url", "")

	input.Sort = app.readString(qs, "sort", "id")

	if !v.Valid() {
		app.errorResponse(w, r, 400, "some error")
		return
	}

	urls, err := app.Storage.Urls.GetAll(input.OriginalUrl, input.ShortUrl)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envelope{"urls": urls}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) ShowUrlHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	url, err := app.Storage.Urls.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJson(w, http.StatusOK, envelope{"url": url}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) UpdateUrlHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}
	url, err := app.Storage.Urls.Get(id)
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
		OriginalUrl *string `json:"original_url"`
		ShortUrl    *string `json:"short_url"`
	}

	err = app.readJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.OriginalUrl != nil {
		url.OriginalUrl = *input.OriginalUrl
	}

	if input.ShortUrl != nil {
		url.ShortUrl = *input.ShortUrl
	}

	err = app.Storage.Urls.Update(url)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envelope{"url": url}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) DeleteUrlHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.Storage.Urls.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJson(w, http.StatusOK, envelope{"message": "success.url is deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
