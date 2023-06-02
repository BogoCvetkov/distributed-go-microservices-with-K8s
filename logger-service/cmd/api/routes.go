package main

import (
	"logger-service/cmd/api/config"
	"logger-service/cmd/api/handlers"
)

func initRoutes(app *config.AppConfig) {

	mux := app.Router

	mux.Post("/create-log", handlers.LogEvent(app))
	mux.Get("/logs", handlers.GetLogs(app))
	mux.Get("/logs/{id}", handlers.GetLog(app))
	mux.Put("/logs/{id}", handlers.UpdateLog(app))
}
