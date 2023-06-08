package config

import (
	data "auth-service/cmd/api/models"
	"net/rpc"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5"
)

type AppConfig struct {
	Router          *chi.Mux
	InitMiddlewares func()
	InitRoutes      func()
	DB              *pgx.Conn
	Models          *data.Models
	RPC             *rpc.Client
}
