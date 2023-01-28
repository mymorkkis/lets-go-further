package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/mymorkkis/lets-go-further-json-api/internal/data"
	"github.com/mymorkkis/lets-go-further-json-api/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if user.Validate(v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// TODO Reduce this time to 1 hour and update user_welcome template
	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {
		data := map[string]any{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}

		if app.mailer.Send(user.Email, "user_welcome.html", data); err == nil {
			properties := map[string]string{"email": user.Email}
			app.logger.PrintInfo("sent welcome email to user", properties)
		} else {
			app.logger.PrintError(err, nil)
		}
	})

	app.serveJSON(w, r, http.StatusAccepted, user, nil)
}
