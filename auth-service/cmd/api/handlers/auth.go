package handlers

import (
	"auth-service/cmd/api/config"
	"auth-service/cmd/api/helpers"
	"fmt"
	"net/http"
)

type jsonResponse struct {
	Err     bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Authenticate(app *config.AppConfig) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var requestPayload struct {
			Email    string `json:"email"`
			Password string `json:"password`
		}

		err := helpers.ParseJSON(w, r, &requestPayload)

		if err != nil {
			helpers.ErrJson(w, fmt.Sprintln(err), http.StatusBadRequest)
			return
		}

		user, err := app.Models.User.GetByEmail(requestPayload.Email)
		if err != nil {
			helpers.ErrJson(w, "Invalid Credentials", http.StatusUnauthorized)
			return
		}

		valid, err := user.PasswordMatches(requestPayload.Password)
		if !valid || err != nil {
			helpers.ErrJson(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		payload := helpers.JsonResponse{
			Message: "Authentication Success",
			Data:    user,
		}

		helpers.WriteJSON(w, &payload)
	}
}
