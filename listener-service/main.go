package main

import (
	"fmt"
	"log"
	"os"
	"time"

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
	log.Println("Connected to RabbitMQ")

	// listen for messages
	// RabbitMQ follows a push approach

	// create consumer

	// watch the queue and consume events

}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue 'till RabbitMQ is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@localhost")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready ...")
			counts++
		} else {
			connection = c
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
