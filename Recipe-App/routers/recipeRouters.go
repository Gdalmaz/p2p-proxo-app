package routers

import (
	"recipe/controllers"
	"recipe/middleware"

	"github.com/gofiber/fiber/v2"
)

func RecipeControllers(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	recipe := v1.Group("/recipe")

	recipe.Post("/create-post", middleware.TokenControl(), controllers.AddRecipe)
	recipe.Put("/update-post/:id", middleware.TokenControl(), controllers.UpdateRecipe)
	recipe.Delete("/delete-post/:id", middleware.TokenControl(), controllers.DeleteRecipe)
}
