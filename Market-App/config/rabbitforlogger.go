package config

import (

	"log"

	"os"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

var (
	RabbitMQConn *amqp.Connection
)

// ConnectRabbitLogger connects to RabbitMQ and declares the queue
func ConnectRabbitLogger() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf(".env yüklenmedi: %v", err)
	}

	rabbitmq := os.Getenv("RABBITMQ_URL")
	if rabbitmq == "" {
		log.Fatal("RABBITMQCONN environment variable not set")
	}
	log.Println(rabbitmq)
	conn, err := amqp.Dial(rabbitmq)
	if err != nil {
		log.Fatalf("RabbitMQ'ya bağlanılamadı: %v", err)
	}
	RabbitMQConn = conn

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("RabbitMQ kanalı oluşturulamadı: %v", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"logger",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Kuyruk oluşturulamadı: %v", err)
	}
}

// SendLog publishes a log message to RabbitMQ

