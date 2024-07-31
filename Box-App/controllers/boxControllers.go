package controllers

import (
	"box/database"
	"box/helpers"
	"box/models"

	"github.com/gofiber/fiber/v2"
)

func CreateBoxFirstShop(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "C-B-F-S-1"})
	}
	if user.IsActive == false {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "C-B-F-S-2"})
	}
	var box models.Box
	err := c.BodyParser(&box)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "C-B-F-S-3"})
	}
	box.UserID = user.ID
	boxControl, _ := helpers.BoxCreateControl(user.ID)
	if boxControl == true {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "C-B-F-S-4"})
	}
	if len(box.Phone) < 9 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "C-B-F-S-5"})
	}

	if len(box.Adress) == 10 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "C-B-F-S-6"})
	}

	err = database.DB.Db.Create(&box).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "C-B-F-S-7"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success"})
}

func UpdateBoxInfo(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-B-I-1"})
	}
	if user.IsActive == false {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "U-B-I-2"})
	}
	var box models.Box
	err := c.BodyParser(&box)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-B-I-3"})
	}
	err = database.DB.Db.Where("user_id=?", user.ID).Updates(&box).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-B-I-4"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success"})
}

func AddProductToBox(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "A-P-T-B-1"})
	}
	id := c.Params("id")
	var box models.Box
	err := database.DB.Db.Where("user_id=?", user.ID).Find(&box).Error
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "A-P-T-B-2"})
	}
	var producttobox models.AddProduct
	err = c.BodyParser(&producttobox)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "A-P-T-B-3", "data": err.Error()})
	}
	producttobox.BoxID = box.ID
	producttobox.UserID = user.ID

	var product models.Product

	err = database.DB.Db.Where("id=?", id).Find(&product).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "A-P-T-B-4"})
	}
	producttobox.ProductID = product.ID
	err = database.DB.Db.Create(&producttobox).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "A-P-T-B-5", "data": err.Error()})
	}
	product.ProductAmount = product.ProductAmount - producttobox.ProductAmount

	err = database.DB.Db.Save(&product).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "A-P-T-B-6"})
	}
	producttobox.Queue++
	box.TotalCost = box.TotalCost + product.ProductPrice
	err = database.DB.Db.Save(&box).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "A-P-T-B-7"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success"})
}

func UpdateBoxProduct(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-B-P-1"})
	}
	var producttobox models.AddProduct
	err := c.BodyParser(&producttobox)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-B-P-2"})
	}
	id := c.Params("id")
	err = database.DB.Db.Where("product_id=? and user_id=?", id, user.ID).Updates(&producttobox).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-B-P-3"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success"})

}

func DeleteProductToBox(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "D-P-T-B-1"})
	}
	var producttobox models.AddProduct
	id := c.Params("id")
	err := database.DB.Db.Where("product_id=? and user_id=?", id, user.ID).Delete(&producttobox).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "D-P-T-B-2"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success"})
}






