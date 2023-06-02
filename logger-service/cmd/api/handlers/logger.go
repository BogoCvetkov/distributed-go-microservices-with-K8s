package handlers

import (
	"fmt"
	"logger-service/cmd/api/config"
	"logger-service/cmd/api/helpers"
	data "logger-service/cmd/api/models"
	"net/http"

	"github.com/go-chi/chi"
)

type jsonResponse struct {
	Err     bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func LogEvent(app *config.AppConfig) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var requestPayload struct {
			EventName string `json:"event_name"`
			Data      string `json:"data`
		}

		err := helpers.ParseJSON(w, r, &requestPayload)

		if err != nil {
			helpers.ErrJson(w, fmt.Sprintln(err), http.StatusBadRequest)
			return
		}

		data := data.LogEntry{
			Name: requestPayload.EventName,
			Data: requestPayload.Data,
		}

		err = app.Models.LogEntry.Insert(&data)

		if err != nil {
			fmt.Println(err)
			helpers.ErrJson(w, "Failed to import log record into DB", http.StatusBadRequest)
			return
		}

		sendResponse(w, &requestPayload)
	}
}

func GetLogs(app *config.AppConfig) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		docs, err := app.Models.LogEntry.GetAll()

		if err != nil {
			fmt.Println(err)
			helpers.ErrJson(w, "Failed to get log records from DB", http.StatusBadRequest)
			return
		}

		sendResponse(w, &docs)

	}
}

func GetLog(app *config.AppConfig) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		logID := chi.URLParam(r, "id")

		doc, err := app.Models.LogEntry.GetOne(logID)

		if err != nil {
			fmt.Println(err)
			helpers.ErrJson(w, fmt.Sprintln(err), http.StatusBadRequest)
			return
		}

		sendResponse(w, &doc)
	}
}

func UpdateLog(app *config.AppConfig) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var requestPayload struct {
			EventName string `json:"event_name"`
			Data      string `json:"data`
		}

		err := helpers.ParseJSON(w, r, &requestPayload)

		if err != nil {
			helpers.ErrJson(w, fmt.Sprintln(err), http.StatusBadRequest)
			return
		}

		data := data.LogEntryUpdate{}

		if requestPayload.EventName != "" {
			data.Name = requestPayload.EventName
		}
		if requestPayload.Data != "" {
			data.Data = requestPayload.Data
		}

		logID := chi.URLParam(r, "id")

		res, err := app.Models.LogEntry.Update(logID, &data)

		if err != nil {
			fmt.Println(err)
			helpers.ErrJson(w, fmt.Sprintln(err), http.StatusBadRequest)
			return
		}

		sendResponse(w, &res)

	}
}

func sendResponse(w http.ResponseWriter, data any) {
	payload := helpers.JsonResponse{
		Message: "Success",
		Data:    data,
	}

	helpers.WriteJSON(w, &payload)
}
