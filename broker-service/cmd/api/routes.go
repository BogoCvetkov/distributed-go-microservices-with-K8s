package main

import (
	"broker-service/cmd/api/config"
	"broker-service/cmd/api/handlers"
	"net/http"
)

func initRoutes(app *config.AppConfig) {

	mux := app.Router

	mux.Post("/", handlers.BrokerMain)
	mux.Post("/route", handlers.RouteRequest(app))
	mux.Get("/hidden", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This route is hidden and should not be exposed outside of the k8s kluster"))
	})
}
