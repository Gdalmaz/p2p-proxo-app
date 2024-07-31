package routers

import (
	"box/controllers"
	"box/middleware"

	"github.com/gofiber/fiber/v2"
)

func BoxRouters(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	box := v1.Group("/box")

	box.Post("/add-product-to-box", middleware.TokenControl(),controllers.AddProductToBox)
	box.Delete("/delete-product-to-box", middleware.TokenControl(),controllers.DeleteProductToBox)
}
