package middleware

import (
	"box/config"
	"box/models"
	"encoding/json"
	"log"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/streadway/amqp"
)


func SendLog(logMessage models.LogMessage) error {
	ch, err := config.RabbitMQConn.Channel()
	if err != nil {
		log.Println("RabbitMQ kanal oluşturulamadı:", err)
		return err
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
		log.Println("Kuyruk oluşturulamadı:", err)
		return err
	}

	messageBody, err := json.Marshal(logMessage)
	if err != nil {
		log.Println("JSON marshalling hatası:", err)
		return err
	}
	log.Println(string(messageBody))
	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        messageBody,
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		message,
	)
	if err != nil {
		log.Println("Mesaj gönderim hatası:", err)
		return err
	}
	log.Println("loggera döndü")
	return nil
}


func LogMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		logMessage := models.LogMessage{
			Timestamp: time.Now().Format(time.RFC3339),
			Status:    c.Response().StatusCode(),
			Latency:   time.Since(start).String(),
			Method:    c.Method(),
			Path:      c.Path(),
		}
		go func() {
			if err := SendLog(logMessage); err != nil {
				log.Println("neyi başaramadın ağğmk")

			}

		}()
		return err
	}

}

