package controllers

import (
	"box/database"
	"box/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AddProductToBox(c *fiber.Ctx) error {
	// Kullanıcı doğrulaması
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized access",
		})
	}

	var req models.AddProduct

	// İstek gövdesini JSON'dan al
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Kullanıcıya ait geçerli bir sepet var mı kontrol et
	var box models.Box
	result := database.DB.Db.Where("user_id = ?", user.ID).First(&box)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Eğer sepet bulunamazsa, yeni bir sepet oluştur
			box = models.Box{
				UserID: user.ID,
			}
			if err := database.DB.Db.Create(&box).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to create a new box",
				})
			}
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve box",
			})
		}
	}

	// Ürünü sepete ekle
	addProduct := models.AddProduct{
		BoxID:         box.ID, // Kullanıcının sepet ID'sini kullan
		ProductID:     req.ProductID,
		ProductAmount: req.ProductAmount,
	}
	var product models.Product
	err := database.DB.Db.Where("id=?", addProduct.ProductID).Find(&product).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "not found product"})
	}
	product.ProductAmount = product.ProductAmount - addProduct.ProductAmount

	// Veritabanına ekle
	if err := database.DB.Db.Create(&addProduct).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add product",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Product added to cart successfully",
		"data":    addProduct,
	})
}

func DeleteProductToBox(c *fiber.Ctx) error {
	return nil
}

//senaryoda her kullanıcının kendine ait bir alışveriş sepeti olucak ve bu alışveriş sepetine bir firmaya ait olmak üzere rastgele ürünler ekleyecek
//eğer farklı bir firmadan bir ürün eklerse sepet otomatik olarak boşaltılıcak
//kullanıcı istediği ürünü istediği miktarda sepete ekleyebilecek ve istediği miktarda sepetten silebilicek .
