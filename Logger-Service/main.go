package main

import (
	"logger/config"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.ConnectRabbitMQ()
	config.ConnectElastic()
	go config.ConsumeRabbit()

	app := fiber.New()
	app.Listen(":9096")
}
