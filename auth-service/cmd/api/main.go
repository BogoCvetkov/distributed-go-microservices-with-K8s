package main

import (
	"auth-service/cmd/api/config"
	data "auth-service/cmd/api/models"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5"
)

const PORT = 3001

func main() {

	// Connect to DB
	conn := connDB()

	// Define app
	app := config.AppConfig{
		Router: chi.NewRouter(),
		DB:     conn,
		Models: data.New(conn),
	}

	// Initialize Middlewares & Routes
	initMiddlewares(&app)
	initRoutes(&app)

	fmt.Printf("Starting auth-service in port %d \n", PORT)

	// http-server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", PORT),
		Handler: app.Router,
	}

	err := srv.ListenAndServe()

	if err != nil {
		fmt.Println("Failed to start auth-service")
		log.Panic(err)
	}
}

func connDB() *pgx.Conn {

	var retries int

	for {
		conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
		if err == nil {
			fmt.Println("Connected to DB")
			return conn
		}

		retries++

		if retries > 10 {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}

		wait := (time.Second / 2) * time.Duration(retries)
		time.Sleep(wait)
		fmt.Printf("DB not ready, retrying after %v \n", wait)
	}

}
