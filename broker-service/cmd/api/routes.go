package main

import (
	"broker-service/cmd/api/handlers"
)

func initRoutes(app *AppConfig) {

	mux := app.router

	mux.Post("/", handlers.BrokerMain)
}
