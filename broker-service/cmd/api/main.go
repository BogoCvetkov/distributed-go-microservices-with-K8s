package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

const PORT = 3000

func main() {

	// Define app
	app := AppConfig{
		router: chi.NewRouter(),
	}

	// Initialize Middlewares & Routes
	initMiddlewares(&app)
	initRoutes(&app)

	fmt.Printf("Starting broker-service in port %d \n", PORT)

	// http-server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", PORT),
		Handler: app.router,
	}

	err := srv.ListenAndServe()

	if err != nil {
		fmt.Println("Failed to start broker-service")
		log.Panic(err)
	}
}
