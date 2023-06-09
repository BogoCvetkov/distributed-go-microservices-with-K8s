package config

import (
	email_proto "broker-service/cmd/api/email_proto"
	"broker-service/cmd/api/types"
	"encoding/json"

	"github.com/go-chi/chi"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AppConfig struct {
	Router          *chi.Mux
	InitMiddlewares func()
	InitRoutes      func()
	RabbitConn      *amqp.Connection
	RabbitChannel   *amqp.Channel
	RabbitQueue     *amqp.Queue
	GClient         email_proto.EmailServiceClient
}

func (app AppConfig) SendToQueue(payload *types.RabbitPayload) error {

	payloadBytes, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	// Publish the payload as a message to the exchange
	err = app.RabbitChannel.Publish(
		"micro_exchange", // exchange name
		"message",        // routing key (empty for direct exchange)
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payloadBytes,
		},
	)

	if err != nil {
		return err
	}

	return nil
}
