package controllers

import (
	"errors"
	"recipe/config"
	"recipe/database"
	"recipe/helpers"
	"recipe/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetRecipe(c *fiber.Ctx) error {
	var click models.Popularity
	var user models.User
	var recipe models.Recipe
	id := c.Params("id")
	err := database.DB.Db.Where("id=?", id).First(&recipe).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "G-R-1"})
	}
	err = database.DB.Db.Where("id=?", recipe.UserID).First(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "G-R-2"})
	}
	err = database.DB.Db.Where("food_id = ? AND user_id = ?", recipe.ID, user.ID).First(&click).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			click.UserID = user.ID
			click.FoodID = recipe.ID
			click.ClickNumber = 1 //First Click
		} else {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "G-R-3"})
		}
	} else {
		click.ClickNumber++
		if click.ClickNumber > 10 && click.ClickNumber%10 == 0 {
			err = config.SaveClickCountToRedis(uint(recipe.ID), click.ClickNumber)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"status": "error", "message": "G-R-4", "data": err.Error()})
			}
			number := click.ClickNumber
			text := "Your recipe is popular now . Your Click Number : \t" + strconv.Itoa(number)
			err = config.RabbitMqPublish([]byte(text), user.Mail)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"status": "error", "message": "G-R-5"})
			}
			err = config.RabbitMqConsume(user.Mail, user.Mail)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"status": "error", "message": "G-R-6"})
			}
		}
	}
	err = database.DB.Db.Save(&click).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "G-R-7"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success", "data": recipe})
}

func GetPopularRecipe(c *fiber.Ctx) error {
	clickcount, err := config.GetAllClickCountsFromRedis()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "G-P-R-1"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success", "data": clickcount})
}

func GetAllRecipe(c *fiber.Ctx) error {
	var recipes []models.Recipe
	err := database.DB.Db.Preload("User").Order("id desc").Find(&recipes).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "G-A-R-1"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success", "data": recipes})
}

func GetUserAllRecipe(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "G-U-A-R-1"})
	}
	Isactive, _ := helpers.CheckVerifyUser(user.ID)
	if Isactive == false {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "G-U-A-R-2"})
	}
	id := c.Params("id")
	var recipe []models.Recipe
	err := database.DB.Db.Preload("User").Where("user_id=?", id).Find(&recipe).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "G-U-A-R-3"})
	}
	
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success", "data": recipe})
}
