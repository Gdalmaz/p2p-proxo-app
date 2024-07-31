package controllers

import (
	"recipe/config"
	"recipe/database"
	"recipe/helpers"
	"recipe/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AddRecipe(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "A-R-1"})
	}
	Isactive, _ := helpers.CheckVerifyUser(user.ID)
	if Isactive == false {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "A-R-2"})
	}
	var recipe models.Recipe
	err := c.BodyParser(&recipe)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "A-R-3"})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "A-R-4"})
	}

	fileBytes, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "A-R-5"})
	}

	defer fileBytes.Close()

	imageBytes := make([]byte, file.Size)
	_, err = fileBytes.Read(imageBytes)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "A-R-6"})
	}

	id, url, err := config.CloudConnect(imageBytes)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "A-R-7"})
	}
	recipe.Image = id
	recipe.ImageUrl = url

	recipe.UserID = user.ID
	foodname := c.FormValue("foodname")
	materials := c.FormValue("materials")
	eatcapacityStr := c.FormValue("eatcapacity")
	description := c.FormValue("description")
	guesspriceStr := c.FormValue("guessprice")
	preparationtimeStr := c.FormValue("preparationtime")

	if len(foodname) != 0 {
		recipe.FoodName = foodname
	}

	if len(materials) != 0 {
		recipe.Materials = materials
	}

	if eatcapacityStr != "" {
		eatcapacity, err := strconv.Atoi(eatcapacityStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "A-R-8"})
		}
		if eatcapacity > 0 {
			recipe.EatCapacity = eatcapacity
		} else {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "A-R-9"})
		}
	}

	if len(description) != 0 {
		recipe.Description = description
	}

	if guesspriceStr != "" {
		guessprice, err := strconv.Atoi(guesspriceStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "A-R-10"})
		}
		if guessprice > 0 {
			recipe.GuessPrice = guessprice
		} else {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "A-R-11"})
		}
	}
	if preparationtimeStr != "" {
		preparationtime, err := strconv.Atoi(preparationtimeStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "A-R-12"})
		}
		if preparationtime > 0 {
			recipe.PrepeareTime = preparationtime
		} else {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "A-R-13"})
		}
	}
	err = database.DB.Db.Create(&recipe).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "A-R-14"})
	}
	text := "Your recipe posted on our website successfully . Thanks ... \n\n\t\t Posted Time : \t" + time.Now().Format("2006-01-02 15:04:05")
	err = config.RabbitMqPublish([]byte(text), user.Mail)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "A-R-15"})
	}
	err = config.RabbitMqConsume(user.Mail, user.Mail)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "A-R-16"})
	}

	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success"})
}

func UpdateRecipe(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "U-R-1"})
	}
	Isactive, _ := helpers.CheckVerifyUser(user.ID)
	if !Isactive {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "U-R-2"})
	}

	recipeID := c.Params("id")
	var recipe models.Recipe
	if err := database.DB.Db.First(&recipe, "id = ? AND user_id = ?", recipeID, user.ID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "U-R-3"})
	}

	// Tarifi güncellemek için gelen form değerlerini mevcut tarifin üzerine yazıyoruz
	if foodname := c.FormValue("foodname"); len(foodname) != 0 {
		recipe.FoodName = foodname
	}

	if materials := c.FormValue("materials"); len(materials) != 0 {
		recipe.Materials = materials
	}

	if eatcapacityStr := c.FormValue("eatcapacity"); eatcapacityStr != "" {
		eatcapacity, err := strconv.Atoi(eatcapacityStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "U-R-8"})
		}
		if eatcapacity > 0 {
			recipe.EatCapacity = eatcapacity
		} else {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "U-R-9"})
		}
	}

	if description := c.FormValue("description"); len(description) != 0 {
		recipe.Description = description
	}

	if guesspriceStr := c.FormValue("guessprice"); guesspriceStr != "" {
		guessprice, err := strconv.Atoi(guesspriceStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "U-R-10"})
		}
		if guessprice > 0 {
			recipe.GuessPrice = guessprice
		} else {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "U-R-11"})
		}
	}

	if preparationtimeStr := c.FormValue("preparationtime"); preparationtimeStr != "" {
		preparationtime, err := strconv.Atoi(preparationtimeStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "U-R-12"})
		}
		if preparationtime > 0 {
			recipe.PrepeareTime = preparationtime
		} else {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "U-R-13"})
		}
	}

	// Yeni bir resim yüklenmişse, resmi işle ve güncelle
	file, err := c.FormFile("image")
	if err == nil {
		fileBytes, err := file.Open()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-R-5"})
		}
		defer fileBytes.Close()

		imageBytes := make([]byte, file.Size)
		_, err = fileBytes.Read(imageBytes)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-R-6"})
		}

		id, url, err := config.CloudConnect(imageBytes)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-R-7"})
		}
		recipe.Image = id
		recipe.ImageUrl = url
	}

	err = database.DB.Db.Save(&recipe).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-R-14"})
	}

	text := "Your recipe updated on our website successfully. Thanks ... \n\n\t\t Updated Time : \t" + time.Now().Format("2006-01-02 15:04:05")
	err = config.RabbitMqPublish([]byte(text), user.Mail)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "U-R-15"})
	}

	err = config.RabbitMqConsume(user.Mail, user.Mail)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "U-R-16"})
	}

	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success"})
}

func DeleteRecipe(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "D-R-1"})
	}
	var recipe models.Recipe
	id := c.Params("id")
	err := database.DB.Db.Where("id=? and user_id=?", id, user.ID).Delete(&recipe).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "D-R-2"})
	}
	err = config.DeleteClickCountFromRedis(uint(recipe.ID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error deleting click count from Redis"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success"})
}
