package main

import (
	"broker-service/cmd/api/config"
	"broker-service/cmd/api/handlers"
)

func initRoutes(app *config.AppConfig) {

	mux := app.Router

	mux.Post("/", handlers.BrokerMain)
	mux.Post("/route", handlers.RouteRequest(app))
}
