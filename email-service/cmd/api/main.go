package main

import (
	"email-service/cmd/api/config"
	"email-service/cmd/api/email_proto"
	mailer "email-service/cmd/api/service"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
	"google.golang.org/grpc"
)

const PORT = 3003
const gRPC_PORT = 9000

func main() {

	// Define app
	app := config.AppConfig{
		Router: chi.NewRouter(),
		Mailer: initMailer(),
	}

	// Initialize Middlewares & Routes
	initMiddlewares(&app)
	initRoutes(&app)

	// initialize gRPC Server
	go initgrpcServer(initMailer())

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

func initgrpcServer(m *mailer.Mail) {

	fmt.Println(fmt.Sprintf("Startin gRPC server on PORT ---> %d", gRPC_PORT))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", gRPC_PORT))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	server := EmailServer{
		Mailer: m,
	}

	email_proto.RegisterEmailServiceServer(s, &server)

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
