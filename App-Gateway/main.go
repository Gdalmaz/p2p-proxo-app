package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:  "*",
		AllowMethods:  "GET,POST,PUT,DELETE",
		AllowHeaders:  "*",
		ExposeHeaders: "*",
		// AllowCredentials: true,
	}))

	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	//AUTH YÖNLENDİRMESİ
	app.Post("/api/v1/user/login", forwardToService("MULTI_AUTH_HOST"))
	app.Post("/api/v1/user/verify-sign-up", forwardToService("MULTI_AUTH_HOST"))
	app.Post("/api/v1/user/signup", forwardToService("MULTI_AUTH_HOST"))
	app.Get("/api/v1/user/logout", forwardToService("MULTI_AUTH_HOST"))
	app.Post("/api/v1/user/delete-account-send-mail", forwardToService("MULTI_AUTH_HOST"))
	app.Delete("/api/v1/user/delete-account-verify-mail", forwardToService("MULTI_AUTH_HOST"))
	app.Put("/api/v1/user/update-password", forwardToService("MULTI_AUTH_HOST"))
	app.Put("/api/v1/user/update-account", forwardToService("MULTI_AUTH_HOST"))
	app.Post("/api/v1/user/forgot-password-help-mail", forwardToService("MULTI_AUTH_HOST"))
	app.Put("/api/v1/user/forgot-password-verify-and-reset", forwardToService("MULTI_AUTH_HOST"))
	//RECİPE YÖNLENDİRMESİ
	app.Post("/api/v1/recipe/create-post", forwardToService("MULTI_RECIPE_HOST"))
	app.Put("/api/v1/recipe/update-post/:id", forwardToService("MULTI_RECIPE_HOST"))
	app.Delete("/api/v1/recipe/delete-post/:id", forwardToService("MULTI_RECIPE_HOST"))
	app.Get("/api/v1/seerec/get-recipe/:id", forwardToService("MULTI_RECIPE_HOST"))
	app.Get("/api/v1/seerec/get-popular-recipe", forwardToService("MULTI_RECIPE_HOST"))
	app.Get("/api/v1/seerec/get-all-recipe", forwardToService("MULTI_RECIPE_HOST"))
	app.Get("/api/v1/seerec/get-user-recipe/:id", forwardToService("MULTI_RECIPE_HOST"))
	//MARKET YÖNLENDİRMESİ
	app.Post("/api/v1/corporation/create-corparation", forwardToService("MULTI_MARKET_HOST"))
	app.Put("/api/v1/corparation/update-corparation", forwardToService("MULTI_MARKET_HOST"))
	app.Put("/api/v1/corparation/verify-update", forwardToService("MULTI_MARKET_HOST"))
	app.Get("/api/v1/corparation/delete-corparation", forwardToService("MULTI_MARKET_HOST"))
	app.Delete("/api/v1/corparation/verify-delete-corparation", forwardToService("MULTI_MARKET_HOST"))
	//PRODUCT YÖNLENDİRMESİ
	app.Post("/api/v1/product/add-product", forwardToService("product/add-product"))
	app.Put("/api/v1/product/update-product/:id", forwardToService("product/update-product/:id"))
	app.Delete("/api/v1/product/delete-product/:id", forwardToService("product/delete-product/:id"))
	//BOX YÖNLENDİRMESİ
	app.Post("/api/v1/box/create-box", forwardToService("box/create-box"))
	app.Put("/api/v1/box/update-box-info", forwardToService("box/update-box-info"))
	app.Post("/api/v1/box/add-product-to-box/:id", forwardToService("box/add-product-to-box/:id"))
	app.Put("/api/v1/box/update-to-product-in-box/:id", forwardToService("box/update-to-product-in-box/:id"))
	app.Delete("api/v1/box/delete-product-to-box/:id", forwardToService("box/delete-product-to-box/:id"))

	log.Fatal(app.Listen(":9095"))
}

func forwardToService(service string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		serviceURL := os.Getenv(service)
		if serviceURL == "" {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("%s is not set", service),
			})
		}

		targetURL := serviceURL + c.Path()
		fmt.Println("Target URL:", targetURL) // Debug: Target URL'yi logla

		body := bytes.NewReader(c.Body())
		req, err := http.NewRequest(c.Method(), targetURL, body)
		if err != nil {
			fmt.Println("NewRequest Error:", err) // Debug: NewRequest hatasını logla
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Headers'ı logla
		fmt.Println("Request Headers:", req.Header)

		req.Header = c.GetReqHeaders()
		client := http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Client Do Error:", err) // Debug: client.Do hatasını logla
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("ReadAll Error:", err) // Debug: ReadAll hatasını logla
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		c.Status(resp.StatusCode)

		for key, value := range resp.Header {
			c.Set(key, value[0])
		}

		return c.Send(bodyBytes)
	}
}
