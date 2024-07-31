package main

import (
	"product/config"
	"product/database"
	"product/routers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoggerRabbit()
	defer config.LoggerRabbit().Close
	database.Connect()
	app := fiber.New()
	routers.ProductRouter(app)
	app.Listen(":9093")
}
