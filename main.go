package main

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

func main() {
	amqpURL := os.Getenv("AMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://admin:admin@192.168.99.100:5672"
	}

	connection, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatal("Error connecting to amqp -", err.Error())
	}

	// Open a channel to publish messages
	channel, err := connection.Channel()
	if err != nil {
		log.Println("Error creating a amqp channel - ", err.Error())
		return
	}

	_, err = channel.QueueDeclare("new_queue", true, false, false, false, nil)
	if err != nil {
		log.Println("Error declaring a queue - ", err.Error())
		return
	}

	// Bind to any changes in the queue. Hence '#'
	err = channel.QueueBind("new_queue", "#", "events", false, nil)
	if err != nil {
		log.Println("Error binding to new_queue -", err.Error())
		return
	}

	// Empty consumer parameter would mean auto generate it. Here we are the consumer process. Ack = false means we need to acknowledge it after processing else it is auto
	// acknowledged. Not an exclusive queue, other consumers can consume as well
	msgs, err := channel.Consume("new_queue", "", false, false, false, false, nil)
	if err != nil {
		log.Fatal("Error consuming the new_queue - ", err.Error())
	}

	// Read continuously
	for msg := range msgs {
		log.Println("MEssages received from the queue:", msg)
		msg.Ack(false)
	}

}
