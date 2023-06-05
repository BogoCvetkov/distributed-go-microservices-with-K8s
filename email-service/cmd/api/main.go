package main

import (
	"email-service/cmd/api/config"
	mailer "email-service/cmd/api/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
)

const PORT = 3003

func main() {

	// Define app
	app := config.AppConfig{
		Router: chi.NewRouter(),
		Mailer: initMailer(),
	}

	// Initialize Middlewares & Routes
	initMiddlewares(&app)
	initRoutes(&app)

	fmt.Printf("Starting email-service in port %d \n", PORT)

	// http-server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", PORT),
		Handler: app.Router,
	}

	err := srv.ListenAndServe()

	if err != nil {
		fmt.Println("Failed to start email-service")
		log.Panic(err)
	}

}

func initMailer() *mailer.Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))

	mailer := mailer.Mail{
		FromAddress: os.Getenv("FROM_ADDR"),
		FromName:    os.Getenv("FROM_NAME"),
		Encryption:  os.Getenv("ENCRYTION"),
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
	}

	return &mailer
}
