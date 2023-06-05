package handlers

import (
	"email-service/cmd/api/config"
	"email-service/cmd/api/helpers"
	mailer "email-service/cmd/api/service"
	"fmt"
	"net/http"
)

func SendEmail(app *config.AppConfig) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var requestPayload struct {
			To      string `json:"to"`
			Subject string `json:"subject`
			Message string `json:"message`
		}

		err := helpers.ParseJSON(w, r, &requestPayload)

		if err != nil {
			helpers.ErrJson(w, fmt.Sprintln(err), http.StatusBadRequest)
			return
		}

		msg := mailer.Message{
			To:      requestPayload.To,
			Subject: requestPayload.Subject,
			Data:    requestPayload.Message,
		}

		err = app.Mailer.SendSMTPMessage(msg)

		if err != nil {
			helpers.ErrJson(w, fmt.Sprintln(err), http.StatusBadRequest)
			return
		}

		payload := helpers.JsonResponse{
			Err:     false,
			Message: fmt.Sprintf("Email send to %s", msg.To),
			Data:    make(map[string]any),
		}

		helpers.WriteJSON(w, &payload)

	}
}
