package main

import (
	"auth/database"
	"auth/routers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.Connect()
	app := fiber.New()
	routers.UserRouter(app)
	app.Listen(":9090")
}
