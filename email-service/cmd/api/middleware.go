package main

import (
	"email-service/cmd/api/config"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func initMiddlewares(app *config.AppConfig) {

	mux := app.Router

	// specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// health check
	mux.Use(middleware.Heartbeat("/ping"))
}
