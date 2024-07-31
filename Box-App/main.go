package main

import (
	"box/database"
	"box/routers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.Connect()
	app := fiber.New()
	routers.BoxRouters(app)
	app.Listen(":80")
}
