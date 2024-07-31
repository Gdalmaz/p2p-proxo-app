package main

import (
	"recipe/config"
	"recipe/database"
	"recipe/routers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.Connect()
	config.ConnectRedis()
	app := fiber.New()
	routers.RecipeControllers(app)
	routers.GetRecipe(app)
	app.Listen(":80")
}
