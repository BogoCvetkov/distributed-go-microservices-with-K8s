package main

import (
	"auth-service/cmd/api/config"
	"auth-service/cmd/api/handlers"
)

func initRoutes(app *config.AppConfig) {

	mux := app.Router

	mux.Post("/auth", handlers.Authenticate(app))
}
