package config

import (
	"auth/middleware"
	"auth/models"
	"log"

	"github.com/streadway/amqp"
)

func Commerce() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)

	}
	defer ch.Close()

	var user models.User

	defer conn.Close()
	q, err := ch.QueueDeclare(
		"authQueue", // Token'ların gönderileceği kuyruk
		true,        // Durable
		false,       // Auto-delete
		false,       // Exclusive
		false,       // No-wait
		nil,         // Arguments
	)
	log.Println(q)
	if err != nil {
		panic(err)
	}

	token, _ := middleware.CreateToken(user.Mail)

	err = publishMessage(ch, "authQueue", token)
	if err != nil {
		panic(err)
	}
	log.Println("market mikroservisine token gönderildi...")

}

func publishMessage(ch *amqp.Channel, queueName string, message string) error {
	err := ch.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return err
	}
	return nil
}
