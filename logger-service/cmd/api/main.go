package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/cmd/api/config"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"

	data "logger-service/cmd/api/models"
	rpc_server "logger-service/cmd/api/rpc"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const PORT = 3002
const RPC_PORT = 5000

func main() {

	// Connect to DB
	conn := connDB()

	// Diconnect context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Close connection on service shut-down
	defer func() {
		if err := conn.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Define app
	app := config.AppConfig{
		Router: chi.NewRouter(),
		DB:     conn,
		Models: data.New(conn),
	}

	// Initialize Middlewares & Routes
	initMiddlewares(&app)
	initRoutes(&app)

	fmt.Printf("Starting logger-service in port %d \n", PORT)

	// http-server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", PORT),
		Handler: app.
			Router,
	}

	// Initialize RPC server
	go setupRPC()

	err := srv.ListenAndServe()

	if err != nil {
		fmt.Println("Failed to start logger-service")
		log.Panic(err)
	}

}

func connDB() *mongo.Client {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URL")))

	if err != nil {
		panic(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to MongoDB")

	return client

}

func setupRPC() {
	rs := new(rpc_server.RPCServer)

	err := rpc.Register(rs)

	if err != nil {
		log.Fatal("Failed to register RPC type ", err)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", RPC_PORT))

	if err != nil {
		log.Panic(err)
	}

	l, err := net.Listen("tcp", tcpAddr.String())

	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Starting logger-service RPC in port %d \n", RPC_PORT)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Spawning new conn")
		go rpc.ServeConn(conn)
	}

}
