package config

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type AppConfig struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Queue   *amqp.Queue
}
