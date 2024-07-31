package controllers

import (
	"product/config"
	"product/database"
	"product/helpers"
	"product/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func AddProduct(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)

	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "A-P-1"})
	}

	// Kullanıcının içerde firmasının olup olmadığını kontrol ediyoruz
	var club models.Club
	err := database.DB.Db.Where("user_id = ?", user.ID).Find(&club).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "A-P-2"})
	}

	var product models.Product
	err = c.BodyParser(&product)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "A-P-3"})
	}

	product.UserID = user.ID
	product.ClubID = club.ID

	productname := c.FormValue("productname")
	productpriceStr := c.FormValue("productprice")
	productamountStr := c.FormValue("productamount")
	if len(productname) != 0 {
		product.ProductName = productname
	}
	if productamountStr != "" {
		productamount, err := strconv.Atoi(productamountStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "A-P-4"})
		}
		if productamount > 0 {
			product.ProductAmount = productamount
		} else {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "A-P-5"})
		}
	}
	if productpriceStr != "" {
		productprice, err := strconv.Atoi(productpriceStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "A-P-6"})
		}
		if productprice > 0 {
			product.ProductPrice = productprice
		} else {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "A-P-7"})
		}
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "A-P-8"})
	}

	fileBytes, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "A-P-9"})
	}
	defer fileBytes.Close()

	imageBytes := make([]byte, file.Size)
	_, err = fileBytes.Read(imageBytes)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "A-P-10"})
	}

	id, url, err := config.CloudConnect(imageBytes)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "A-P-11"})
	}

	product.Image = id
	product.ImageUrl = url

	err = database.DB.Db.Create(&product).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "A-P-12"})
	}
	err = database.DB.Db.Preload("User").Preload("Club").First(&product, product.ID).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "A-P-13"})
	}

	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success", "data": product})
}
func UpdateProduct(c *fiber.Ctx) error {
	// User ve Club bilgilerini alıyoruz
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-P-1"})
	}

	var club models.Club
	err := database.DB.Db.Where("user_id=?", user.ID).Find(&club).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-P-2"})
	}

	if user.ID != club.UserID {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-P-3"})
	}

	// Ürünü ID'sine göre alıyoruz
	id := c.Params("id")
	var product models.Product
	err = database.DB.Db.Where("id = ? AND club_id = ?", id, club.ID).First(&product).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "U-P-4"})
	}

	// Gönderilen form verilerini alıyoruz
	productname := c.FormValue("productname")
	productpriceStr := c.FormValue("productprice")
	productamountStr := c.FormValue("productamount")

	if len(productname) != 0 {
		product.ProductName = productname
	}
	if productamountStr != "" {
		productamount, err := strconv.Atoi(productamountStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "U-P-5"})
		}
		if productamount > 0 {
			product.ProductAmount = productamount
		} else {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "U-P-6"})
		}
	}
	if productpriceStr != "" {
		productprice, err := strconv.Atoi(productpriceStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "U-P-7"})
		}
		if productprice > 0 {
			product.ProductPrice = productprice
		} else {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "U-P-8"})
		}
	}

	// Image dosyasını alıyoruz
	file, err := c.FormFile("image")
	if err == nil {
		fileBytes, err := file.Open()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-P-9"})
		}
		defer fileBytes.Close()

		imageBytes := make([]byte, file.Size)
		_, err = fileBytes.Read(imageBytes)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-P-10"})
		}

		id, url, err := config.CloudConnect(imageBytes)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-P-11"})
		}
		product.Image = id
		product.ImageUrl = url
	}

	// Ürünü güncelliyoruz
	err = database.DB.Db.Save(&product).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-P-12"})
	}

	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Product updated successfully", "data": product})
}

func DeleteProduct(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "V-D-P-1"})
	}
	var club models.Club

	err := database.DB.Db.Where("user_id=?", user.ID).Find(&club).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "V-D-P-2"})
	}
	if user.ID != club.UserID {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "V-D-P-3"})
	}
	var product models.Product
	id := c.Params("id")
	var password models.VerifyPass
	err = c.BodyParser(&password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "V-D-P-4"})
	}

	password.UserPassword = user.Password
	password.InputPassword = helpers.HashPass(password.InputPassword)
	if password.UserPassword != password.InputPassword {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "V-D-P-5"})
	}

	err = database.DB.Db.Where("id=?", id).Delete(&product).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "V-D-P-6"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success"})

}
