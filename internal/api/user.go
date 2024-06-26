package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/yantay0/url-shortener/internal/model"
	"github.com/yantay0/url-shortener/internal/storage"
	"github.com/yantay0/url-shortener/internal/validator"
)

func (app *App) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	// Create an anonymous struct to hold the expected data from the request body.
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse the request body into the anonymous struct
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &model.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	// Use the Password.Set() method to generate and store the hashed and plaintext
	// passwords.
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if model.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.Storage.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.Storage.Permissions.AddForUser(user.ID, "shortenings:read")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// After the user record has been created in the database, generate a new activation
	// token for the user.
	token, err := app.Storage.Tokens.New(user.ID, 3*24*time.Hour, storage.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Launch a goroutine which runs an anonymous function that sends the welcome email using
	// the background helper function.
	app.background(func() {
		// Create map to act as a 'holding structure' for the data we send to the weclome email
		// template.
		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}

		// Call the Send() method on our Mailer, passing in the user's email address, name of the
		// template file, and the data map containing the activationToken and the user's ID.
		err = app.Mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			// Importantly, if there is an error sending the email then we log the error
			// instead of raising a server error like before when we handled
			// the email send functionality without a goroutine
			app.Logger.PrintError(err, nil)
		}
	})

	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the plaintext activation token from the request body.
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if model.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.Storage.Users.GetForToken(storage.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	user.Activated = true
	// Save the updated user record in our database, checking for any edit conflicts in
	// the same way that we did for our movie records.
	err = app.Storage.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.Storage.Tokens.DeleteAllForUser(storage.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Send the updated user details to the client in a JSON response.
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
