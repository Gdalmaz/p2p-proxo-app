package config

import (
	"log"

	"github.com/streadway/amqp"
)

var (
	conn *amqp.Connection 
	rabbitMQCh *amqp.Channel
)

func ConnectRabbitMQ() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:15672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
}

func LoggerRabbit() {

	// RabbitMQ'ya bağlan

	defer conn.Close()

	rabbitMQCh, err = conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer rabbitMQCh.Close()

	_, err = rabbitMQCh.QueueDeclare(
		"log_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

}

func sendLogMessage(message string) {
	err := rabbitMQCh.Publish(
		"",
		"log_queue",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		log.Printf("Failed to publish a message: %s", err)
	}
}
