package main

import (
	"broker-service/cmd/api/config"
	"broker-service/cmd/api/helpers"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/go-chi/chi"
)

const PORT = 3000

func main() {

	// Connect to message broker
	conn, err := connectRabbitMQ()

	if err != nil {
		fmt.Println("Failed to connect to RabbitMQ")
		log.Panic(err)
	}

	defer conn.Close()

	// Define app
	app := config.AppConfig{
		Router:     chi.NewRouter(),
		RabbitConn: conn,
	}

	// Initialize Middlewares & Routes
	initMiddlewares(&app)
	initRoutes(&app)

	// Init exchanges and queues
	helpers.PrepareRabbitConn(&app)

	defer app.RabbitChannel.Close()

	fmt.Printf("Starting broker-service in port %d \n", PORT)

	// http-server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", PORT),
		Handler: app.Router,
	}

	err = srv.ListenAndServe()

	if err != nil {
		fmt.Println("Failed to start broker-service")
		log.Panic(err)
	}
}

func connectRabbitMQ() (*amqp.Connection, error) {
	var retries int

	for {
		conn, err := amqp.Dial(os.Getenv("RABBIT_URL"))

		if err == nil {
			fmt.Println("Connected to RabbitMQ")
			return conn, nil
		}

		retries++

		if retries > 10 {
			fmt.Fprintf(os.Stderr, "Unable to connect to RabbitMQ: %v\n", err)
			return nil, err
		}

		wait := (time.Second / 2) * time.Duration(retries)

		time.Sleep(wait)
		fmt.Printf("RabbitMQ not ready, retrying after %v \n", wait)

	}

}
