package routers

import (
	"product/controllers"
	"product/middleware"

	"github.com/gofiber/fiber/v2"
)

func ProductRouter(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	product := v1.Group("/product")

	product.Post("/add-product", middleware.TokenControl(), controllers.AddProduct)
	product.Put("/update-product/:id", middleware.TokenControl(), controllers.UpdateProduct)
	product.Delete("/delete-product/:id", middleware.TokenControl(), controllers.DeleteProduct)

}
