package handlers

import (
	"broker-service/cmd/api/helpers"
	"net/http"
)

type jsonResponse struct {
	Err     bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func BrokerMain(w http.ResponseWriter, r *http.Request) {

	payload := helpers.JsonResponse{
		Message: "Greetings from the broker service",
	}

	helpers.WriteJSON(w, &payload)

}
