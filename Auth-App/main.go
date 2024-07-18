package main

import (
	"auth/Auth-App/database"
	"auth/Auth-App/routers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.Connect()
	app := fiber.New()
	routers.UserRouter(app)
	app.Listen(":80")
}
