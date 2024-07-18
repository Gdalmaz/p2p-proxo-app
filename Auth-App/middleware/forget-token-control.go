package middleware

import (
	"auth/Auth-App/database"
	"auth/Auth-App/models"

	"github.com/gofiber/fiber/v2"
)

func ForgotTokenControl() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.DB.Db
		authorizationHeader := c.Get("Authorization")
		if authorizationHeader == "" || len(authorizationHeader) < 7 || authorizationHeader[:7] != "Bearer " {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Geçersiz veya eksik token",
			})
		}
		token := authorizationHeader[7:]

		var session models.ForgotPassword
		if err := db.Where("token=?", token).First(&session).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Geçersiz token",
				"message": "Oturum bulunamadı veya süresi dolmuş",
			})
		}

	


		// Kullanıcıyı Fiber bağlamına ekleyin
		c.Locals("user", session)

		return c.Next()

	}
}
