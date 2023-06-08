package handlers

import (
	"auth-service/cmd/api/config"
	"auth-service/cmd/api/helpers"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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

		err = logRPC(app, "User Logged", fmt.Sprintf("RPC --> %s %s has logged in at %v", user.FirstName, user.LastName, time.Now()))

		if err != nil {
			fmt.Println(err)
			helpers.ErrJson(w, fmt.Sprintln(err), http.StatusBadRequest)
			return
		}

		payload := helpers.JsonResponse{
			Err:     false,
			Message: "Authentication Success",
			Data:    user,
		}

		helpers.WriteJSON(w, &payload)
	}
}

func logRequest(name string, data string) error {
	var entry struct {
		Name string `json:"event_name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, err := json.Marshal(entry)

	if err != nil {
		return err
	}

	url := "http://logger-service:3002/create-log"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(req)

	if err != nil {
		return err
	}

	return nil
}

func logRPC(app *config.AppConfig, name string, data string) error {

	type RPCNewLog struct {
		EventName string
		Data      string
	}

	entry := RPCNewLog{
		EventName: name,
		Data:      data,
	}

	var reply jsonResponse

	err := app.RPC.Call("RPCServer.LogEntry", entry, &reply)

	fmt.Println(reply)

	if err != nil {
		return err
	}

	return nil
}
