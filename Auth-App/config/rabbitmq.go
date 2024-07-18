package config

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"github.com/streadway/amqp"
)

// RabbitMqPublish publishes a message to a specified queue
func RabbitMqPublish(msg []byte, queueName string) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		queueName, // queue name
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %v", err)
	}

	message := amqp.Publishing{
		ContentType: "text/plain",
		Body:        msg,
	}

	err = ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		message,
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %v", err)
	}

	log.Println("Message successfully published!")
	return nil
}

// RabbitMqConsume consumes messages from a specified queue and sends emails
func RabbitMqConsume(queueName, sender string) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	messages, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Printf("Successfully connected to queue: %s", queueName)

	go func() {
		for message := range messages {
			body := string(message.Body)
			to := []string{queueName}
			if err := SendMail(to, sender, body); err != nil {
				log.Printf("Failed to send email: %v", err)
			} else {
				log.Printf("Email successfully sent to: %s", to[0])
			}
		}
	}()
	return nil
}

// SendMail sends an email using the specified parameters
func SendMail(to []string, sender string, body string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "465"
	smtpUsername := "gokhandalmzz@gmail.com"
	smtpPassword := "bfkq mmiu vuij nrnb"

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)
	msg := "From: " + sender + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Subject: Test Email\n\n" + body

	// Dial TCP connection
	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	})
	if err != nil {
		return fmt.Errorf("TLS connection error: %v", err)
	}
	defer conn.Close()

	// Create a new SMTP client
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return fmt.Errorf("error creating SMTP client: %v", err)
	}
	defer client.Quit()

	// Auth
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP authentication error: %v", err)
	}

	// Set sender address
	if err = client.Mail(smtpUsername); err != nil {
		return fmt.Errorf("sender address error: %v", err)
	}
	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("recipient address error: %v", err)
		}
	}

	// Data
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("data sending error: %v", err)
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("message writing error: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("message closing error: %v", err)
	}

	log.Println("Email successfully sent")
	return nil
}
