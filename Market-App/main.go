package main

import (
	"market/config"
	"market/database"
	"market/middleware"
	"market/routers"

	"github.com/gofiber/fiber/v2"
)

func main() {

	database.Connect()
	config.ConnectRabbitLogger()
	app := fiber.New()
	app.Use(middleware.LogMiddleware())
	routers.CorporationRouters(app)
	app.Listen(":9092")
}
