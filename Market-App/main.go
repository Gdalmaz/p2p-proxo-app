package main

import (
	"market/database"
	"market/routers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.Connect()
	app := fiber.New()
	routers.CorporationRouters(app)
	app.Listen(":80")
}
