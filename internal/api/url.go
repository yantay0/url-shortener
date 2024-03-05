package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/yantay0/url-shortener/internal/model"
	"github.com/yantay0/url-shortener/internal/storage"
)

func (app *App) CreateUrlHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Original_url string `json:"original_url"`
		Short_url    string `json:"short_url"`
	}

	err := app.readJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	url := &model.Url{
		Original_url: input.Original_url,
		Short_url:    input.Short_url,
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
	fmt.Fprint(w, "lsit")
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
		case errors.Is(err, storage.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	var input struct {
		Original_url string `json:"original_url"`
		Short_url    string `json:"short_url"`
	}

	err = app.readJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	url.Original_url = input.Original_url
	url.Short_url = input.Short_url

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
