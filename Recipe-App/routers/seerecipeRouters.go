package routers

import (
	"recipe/controllers"
	"recipe/middleware"

	"github.com/gofiber/fiber/v2"
)

func GetRecipe(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	seerec := v1.Group("/seerec")

	seerec.Get("/get-recipe/:id", controllers.GetRecipe)
	seerec.Get("/get-popular-recipe", controllers.GetPopularRecipe)
	seerec.Get("/get-all-recipe", controllers.GetAllRecipe)
	seerec.Get("/get-user-recipe/:id", middleware.TokenControl(), controllers.GetUserAllRecipe)
}
