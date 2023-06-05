package main

import (
	"email-service/cmd/api/config"
	"email-service/cmd/api/handlers"
)

func initRoutes(app *config.AppConfig) {

	mux := app.Router

	mux.Post("/send", handlers.SendEmail(app))
}
