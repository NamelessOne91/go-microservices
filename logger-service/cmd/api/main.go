package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/NamelessOne91/logger/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	grpcPort = "50001"
)

type Config struct {
	Repo data.Repository
}

func main() {
	// connect to Mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	//close connection
	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Repo: data.NewMongoRepository(mongoClient),
	}

	// Register the RPC server
	err = rpc.Register(
		RPCServer{
			repo: data.NewMongoRepository(mongoClient),
		},
	)
	go app.rpcListen()

	go app.gRPCListen()

	// start web server
	log.Println("Starting service on port", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func (app *Config) rpcListen() error {
	log.Println("Starting RPC server on port", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}

func connectToMongo() (*mongo.Client, error) {
	// create connection options
	clientOpts := options.Client().ApplyURI(mongoURL)
	// read from env in production !
	clientOpts.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	c, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Println("Error connecting to Mongo:", err)
	}
	return c, nil
}
