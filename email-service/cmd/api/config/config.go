package config

import (
	mailer "email-service/cmd/api/service"

	"github.com/go-chi/chi"
)

type AppConfig struct {
	Router          *chi.Mux
	InitMiddlewares func()
	InitRoutes      func()
	Mailer          *mailer.Mail
}
