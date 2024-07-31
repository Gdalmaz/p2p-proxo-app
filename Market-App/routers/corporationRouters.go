package routers

import (
	"market/controllers"
	"market/middleware"

	"github.com/gofiber/fiber/v2"
)

func CorporationRouters(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	corporation := v1.Group("/corporation")

	corporation.Post("/create-corparation", middleware.TokenControl(), controllers.CreateCorporation)
	corporation.Put("/update-corparation", middleware.TokenControl(), controllers.UpdateCorparation)
	corporation.Put("/verify-update", middleware.TokenControl(), controllers.VerifySendUpdate)
	corporation.Get("/delete-corparation", middleware.TokenControl(), controllers.DeleteCorparation)
	corporation.Delete("/verify-delete-corparation", middleware.TokenControl(), controllers.VerifyDeleteCorparation)
}
