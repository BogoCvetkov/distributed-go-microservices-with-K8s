package handlers

import (
	"broker-service/cmd/api/config"
	email_proto "broker-service/cmd/api/email_proto"
	"broker-service/cmd/api/helpers"
	"broker-service/cmd/api/types"
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"
)

func BrokerMain(w http.ResponseWriter, r *http.Request) {

	payload := helpers.JsonResponse{
		Message: "Greetings from the broker service",
	}

	helpers.WriteJSON(w, &payload)

}

// Used to route requests to the different micro-services
func RouteRequest(app *config.AppConfig) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var payload types.RouteRequestBody

		err := helpers.ParseJSON(w, r, &payload)

		if err != nil {
			helpers.ErrJson(w, fmt.Sprintf("Failed parsing JSON - %s", err))
			return
		}

		switch payload.Action {
		case "auth":
			authenticate(w, r, &payload)
			return
		case "log":
			log(w, r, &payload)
			return
		case "send":
			send(w, r, &payload, app)
			return
		case "send-grpc":
			sendgRPC(w, r, &payload, app)
			return
		default:
			helpers.ErrJson(w, "Unrecognized action")
			return
		}
	}

}

func authenticate(w http.ResponseWriter, r *http.Request, data *types.RouteRequestBody) {

	type AuthPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	dataMap, ok := data.Payload.(map[string]any)
	if !ok {
		helpers.ErrJson(w, "payload field not a map")
		return
	}

	payload := AuthPayload{
		Email:    dataMap["email"].(string),
		Password: dataMap["password"].(string),
	}

	reqInfo := types.MethodCallInfo{
		Method:   "POST",
		Endpoint: fmt.Sprintf("%s/auth", os.Getenv("AUTH_ENDPOINT")),
		Body:     payload,
	}

	helpers.CallModule(w, &reqInfo)

}

func log(w http.ResponseWriter, r *http.Request, data *types.RouteRequestBody) {

	type LogPayload struct {
		Name string `json:"event_name"`
		Data string `json:"data"`
	}

	dataMap, ok := data.Payload.(map[string]any)
	if !ok {
		helpers.ErrJson(w, "payload field not a map")
		return
	}

	payload := LogPayload{
		Name: dataMap["event_name"].(string),
		Data: dataMap["data"].(string),
	}
	reqInfo := types.MethodCallInfo{
		Method:   "POST",
		Endpoint: fmt.Sprintf("%s/create-log", os.Getenv("LOGGER_ENDPOINT")),
		Body:     payload,
	}

	helpers.CallModule(w, &reqInfo)
}

func sendgRPC(w http.ResponseWriter, r *http.Request, data *types.RouteRequestBody, app *config.AppConfig) {
	type Data struct {
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	dataMap, ok := data.Payload.(map[string]any)
	if !ok {
		helpers.ErrJson(w, "payload field not a map")
		return
	}

	gRequest := email_proto.EmailRequest{
		Data: &email_proto.EmailData{
			To:      dataMap["to"].(string),
			Subject: dataMap["subject"].(string),
			Message: dataMap["message"].(string),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	gres, err := app.GClient.SendEmail(ctx, &gRequest)

	if err != nil {
		helpers.ErrJson(w, fmt.Sprintln(err))
		return
	}

	res := helpers.JsonResponse{
		Err:     false,
		Message: gres.Message,
	}

	helpers.WriteJSON(w, &res)

}

func send(w http.ResponseWriter, r *http.Request, data *types.RouteRequestBody, app *config.AppConfig) {

	type Data struct {
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	dataMap, ok := data.Payload.(map[string]any)
	if !ok {
		helpers.ErrJson(w, "payload field not a map")
		return
	}

	d := Data{
		To:      dataMap["to"].(string),
		Subject: dataMap["subject"].(string),
		Message: dataMap["message"].(string),
	}
	payload := types.RabbitPayload{
		Endpoint: fmt.Sprintf("%s/send", os.Getenv("EMAIL_ENDPOINT")),
		Method:   "POST",
		Data:     d,
	}

	if err := app.SendToQueue(&payload); err != nil {
		helpers.ErrJson(w, fmt.Sprintln(err))
		return
	}

	res := helpers.JsonResponse{
		Err:     false,
		Message: "Message send to broker",
	}

	helpers.WriteJSON(w, &res)

}
