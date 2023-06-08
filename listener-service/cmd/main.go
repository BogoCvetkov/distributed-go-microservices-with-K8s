package main

import (
	"fmt"
	"listener-service/cmd/config"
	"listener-service/cmd/helpers"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

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
		Conn: conn,
	}

	// Init exchanges and queues
	helpers.PrepareRabbitConn(&app)

	defer app.Channel.Close()

	// Declare blocking channel
	var forever chan struct{}

	// Start listening for messages and consume them
	go helpers.ListenOnQueue(&app)

	fmt.Printf("Started listening for messages on Queue --> %s ", app.Queue.Name)

	log.Printf("To exit press CTRL+C")
	<-forever
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
