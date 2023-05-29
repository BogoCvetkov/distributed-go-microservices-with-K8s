package main

import "github.com/go-chi/chi"

type AppConfig struct {
	router          *chi.Mux
	initMiddlewares func()
	initRoutes      func()
}
