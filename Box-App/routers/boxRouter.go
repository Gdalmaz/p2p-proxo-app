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

	box.Post("/create-box", middleware.TokenControl(), controllers.CreateBoxFirstShop)
	box.Put("/update-box-info", middleware.TokenControl(), controllers.UpdateBoxInfo)
	box.Post("/add-product-to-box/:id", middleware.TokenControl(), controllers.AddProductToBox)
	box.Put("/update-to-product-in-box/:id", middleware.TokenControl(), controllers.UpdateBoxProduct)
	box.Delete("/delete-product-to-box/:id", middleware.TokenControl(), controllers.DeleteProductToBox)
}
