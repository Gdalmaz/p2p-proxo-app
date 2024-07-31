package main

import (
	"box/config"
	"box/database"
	"box/middleware"
	"box/routers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.Connect()
	config.ConnectRabbitLogger()
	app := fiber.New()
	app.Use(middleware.LogMiddleware())
	routers.BoxRouters(app)
	app.Listen(":9094")
}
