package config

import (
	data "logger-service/cmd/api/models"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppConfig struct {
	Router          *chi.Mux
	InitMiddlewares func()
	InitRoutes      func()
	DB              *mongo.Client
	Models          *data.Models
}
