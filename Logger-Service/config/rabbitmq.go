package config

import (
	"encoding/json"
	"log"
	"logger/models"
	"time"

	"os"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

var (
	Elastic *elasticsearch.Client
	Conn    *amqp.Connection
)
//CONNECT RABBİT
func ConnectRabbitMQ() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("hata aldık .env yüklenmedi")
		panic(err)
	}
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		log.Fatal("RABBITMQ_URL environment variable is not set")
	}
	log.Println(rabbitMQURL)
	Conn, err = amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	log.Println("Connected to RabbitMQ:", Conn)
}

//CONSUME RABBİT
func ConsumeRabbit() {
	if Conn == nil {
		log.Fatal("RabbitMQ connection is nil. Make sure to call ConnectRabbitMQ first.")
	}

	log.Println("Starting to consume messages from RabbitMQ")

	ch, err := Conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"logger",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}
	log.Println(msgs)
	go func ()  {
		for d := range msgs {
			logMessage := models.LogMessage{}
				messageBody:= string(d.Body)
				log.Println(messageBody)
				if err := json.Unmarshal([]byte(messageBody),&logMessage) ; err != nil{
					continue
				}
				if err := SendLogToElasticsearch(logMessage); err !=nil{
					log.Println("vahiy gitmedi")
				}

		}
	}()
	for {
		time.Sleep(time.Second)
	}
}

