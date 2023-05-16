package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/NamelessOne91/listener/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// connect to RabbitMQ
	conn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	// listen for messages
	// RabbitMQ follows a push approach
	log.Println("Listening for and consuming RabbitMQ messages ...")

	// create consumer
	consumer, err := event.NewConsumer(conn)
	if err != nil {
		panic(err)
	}
	// watch the queue and consume events
	err = consumer.Listen("log.INFO", "log.WARNING", "log.ERROR")
	if err != nil {
		log.Println(err)
	}

}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue 'till RabbitMQ is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready ...")
			counts++
		} else {
			connection = c
			log.Println("Connected to RabbitMQ")
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		// exponential backoff
		backOff = time.Duration(2<<counts) * time.Second
		log.Println("backing off ...")
		time.Sleep(backOff)
		continue
	}
	return connection, nil
}
