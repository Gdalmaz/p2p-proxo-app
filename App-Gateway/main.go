package main

import (
    "log"
    "os"
    "strings"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/proxy"
    "github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    app := fiber.New()

    app.All("/api/v1/user/*", func(c *fiber.Ctx) error {
        target := os.Getenv("MULTI_AUTH_HOST")
        return proxy.Do(c, target+strings.TrimPrefix(c.OriginalURL(), "/api/v1/user"))
    })

    app.All("/api/v1/recipe/*", func(c *fiber.Ctx) error {
        target := os.Getenv("MULTI_RECIPE_HOST")
        return proxy.Do(c, target+strings.TrimPrefix(c.OriginalURL(), "/api/v1/recipe"))
    })

    app.All("/api/v1/seerec/*", func(c *fiber.Ctx) error {
        target := os.Getenv("MULTI_VISIT_RECIPE_HOST")
        return proxy.Do(c, target+strings.TrimPrefix(c.OriginalURL(), "/api/v1/seerec"))
    })

    app.All("/api/v1/corporation/*", func(c *fiber.Ctx) error {
        target := os.Getenv("MULTI_MARKET_HOST")
        return proxy.Do(c, target+strings.TrimPrefix(c.OriginalURL(), "/api/v1/corporation"))
    })

    app.All("/api/v1/product/*", func(c *fiber.Ctx) error {
        target := os.Getenv("MULTI_PRODUCT_HOST")
        return proxy.Do(c, target+strings.TrimPrefix(c.OriginalURL(), "/api/v1/product"))
    })

    log.Fatal(app.Listen(":8080"))
}