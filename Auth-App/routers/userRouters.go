package routers

import (
	"auth/controllers"
	"auth/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	user := v1.Group("/user")

	user.Post("/signup", controllers.SignUp)
	user.Post("/verify-sign-up", middleware.VerifyTokenControl(), controllers.LastStepSignUp)
	user.Post("/login", controllers.LogIn)
	user.Get("/logout", middleware.TokenControl(), controllers.LogOut)
	user.Post("/delete-account-send-mail", middleware.TokenControl(), controllers.DeleteAccountSendMail)
	user.Delete("/delete-account-verify-mail", middleware.TokenControl(), controllers.DeleteAccountVerifyMail)
	user.Put("/update-password", middleware.TokenControl(), controllers.UpdatePassword)
	user.Put("/update-account", middleware.TokenControl(), controllers.UpdateAccount)
	user.Post("/forgot-password-help-mail", controllers.ForgotPasswordSendCode)
	user.Put("/forgot-password-verify-and-reset", middleware.ForgotTokenControl(), controllers.ForgotPasswordVerifyAndReset)
}
