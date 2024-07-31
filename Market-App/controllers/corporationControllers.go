package controllers

import (
	"market/config"
	"market/database"
	"market/helpers"
	"market/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateCorporation(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		c.Status(404).JSON(fiber.Map{"status": "error", "message": "C-C-1"})
	}
	var club models.Club
	err := c.BodyParser(&club)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "C-C-2"})
	}
	if user.IsActive == false {
		return c.Status(200).JSON(fiber.Map{"status": "error", "message": "C-C-3"})
	}
	club.UserFullName = user.FullName
	club.UserID = user.ID

	clubControl, _ := helpers.ClubControl(user.ID)
	if clubControl == true {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "C-C-4"})
	}
	clubname := c.FormValue("clubname")
	explanation := c.FormValue("explanation")
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "C-C-5"})
	}

	fileBytes, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "C-C-6"})
	}

	defer fileBytes.Close()

	imageBytes := make([]byte, file.Size)
	_, err = fileBytes.Read(imageBytes)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "C-C-7"})
	}

	id, url, err := config.CloudConnect(imageBytes)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "C-C-8"})
	}
	club.Image = id
	club.ImageUrl = url

	if len(clubname) != 0 {
		club.ClubName = clubname
	}
	if len(explanation) != 0 {
		club.Explanation = explanation
	}

	err = database.DB.Db.Create(&club).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "C-C-8"})
	}
	Text := "Your corparation created successfuly.Created Time :" + time.Now().Format("2006-01-02 15:04:05")
	err = config.RabbitMqPublish([]byte(Text), user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "C-C-9"})
	}
	err = config.RabbitMqConsume(user.Mail, user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "C-C-10"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success"})
}

func UpdateCorparation(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-C-1"})
	}
	var club models.Club
	err := c.BodyParser(&club)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-C-2"})
	}
	err = database.DB.Db.Where("user_id = ?", user.ID).First(&club).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "U-C-3"})
	}
	clubname := c.FormValue("clubname")
	explanation := c.FormValue("explanation")
	file, err := c.FormFile("image")

	if len(clubname) != 0 {
		club.ClubName = clubname
	}
	if len(explanation) != 0 {
		club.Explanation = explanation
	}

	if file != nil {
		fileBytes, err := file.Open()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-C-4"})
		}
		defer fileBytes.Close()

		imageBytes := make([]byte, file.Size)
		_, err = fileBytes.Read(imageBytes)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-C-5"})
		}

		id, url, err := config.CloudConnect(imageBytes)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-C-6"})
		}
		club.Image = id
		club.ImageUrl = url
	}

	err = database.DB.Db.Save(&club).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-C-7"})
	}
	Text := "Your corporation updated successfully. Updated Time: " + time.Now().Format("2006-01-02 15:04:05")
	err = config.RabbitMqPublish([]byte(Text), user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-C-8"})
	}

	err = config.RabbitMqConsume(user.Mail, user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-C-9"})
	}

	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Corporation updated successfully"})
}

func VerifySendUpdate(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "V-S-U-1"})
	}
	var systemCode models.VerifyUpdate
	code := helpers.GenerateRandomNumber()
	systemCode.Code = code
	codeStr := strconv.Itoa(code)
	err := config.RabbitMqPublish([]byte(codeStr), user.Mail)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "V-S-U-2"})
	}
	err = config.RabbitMqConsume(user.Mail, user.Mail)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "V-S-U-3"})
	}
	err = database.DB.Db.Create(&systemCode).Error
	if err != nil {
		return c.Status(200).JSON(fiber.Map{"status": "error", "message": "V-S-U-4"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success"})
}

func VerifyPutUpdate(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "V-P-U-1"})
	}
	var verifyUpdate models.VerifyUpdate
	err := database.DB.Db.Where("user_id=?", user.ID).Find(&verifyUpdate).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "V-P-U-2"})
	}
	var code models.Code
	err = c.BodyParser(&code)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "V-P-U-3"})
	}
	code.SystemCode = verifyUpdate.Code
	if code.SystemCode != code.InputCode {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "V-P-U-4"})
	}
	err = database.DB.Db.Where("user_id=?", user.ID).Find(&verifyUpdate).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "V-P-U-5"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success"})
}

func DeleteCorparation(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "D-C-1"})
	}
	var systemCode models.VerifyDelete
	code := helpers.GenerateRandomNumber()
	systemCode.Code = code
	systemCode.UserID = user.ID
	codeStr := strconv.Itoa(code)
	err := config.RabbitMqPublish([]byte(codeStr), user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "D-C-2"})
	}
	err = config.RabbitMqConsume(user.Mail, user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "D-C-3"})
	}

	err = database.DB.Db.Create(&systemCode).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "D-C-4"})
	}

	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success"})
}

func VerifyDeleteCorparation(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "V-D-C-1"})
	}
	var code models.Code
	err := c.BodyParser(&code)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "V-D-C-2"})
	}
	var systemCode models.VerifyDelete
	err = database.DB.Db.Where("user_id=?", user.ID).Find(&systemCode).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "V-D-C-3"})
	}
	code.SystemCode = systemCode.Code
	if code.SystemCode != code.InputCode {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "V-D-C-4"})
	}
	var club models.Club
	deleteID := user.ID

	err = database.DB.Db.Where("user_id=?", user.ID).First(&club).Delete(&club).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "V-D-C-5"})
	}

	err = database.DB.Db.Where("user_id=?", deleteID).Delete(&systemCode).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "V-D-C-6"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success"})
}

//										UPDATE SENARYO

// Tasarlanan senaryoya göre kullanıcı firma bilgilerini update etmek istediği zaman front end tarafından
// tasarlanan update butonuna tıklandıktan sonra mail doğrulaması yapması gerekecek eğer 200 dönerse
// update sayfasına yönlendirilecek ve kullanıcı burada firma bilgilerini başarılı bir şekilde güncelleyebilicek

//										DELETE SENARYO
//Bu kısımdada aynı şekilde update kısmında olduğu gibi bir senaryo vardır

//NOTE Buradaki senaryo daha farklı algoritmalardada yapılabilir mesela kullanıcıya ait gönderilen kodların tutulduğu tabloya isactive adında bir veri koyulur
// ve bunun default hali false olur kullanıcı kodu onayladığında otomaik true olur ve aktif olduğunda silme işlemini gerçekleştirebiliriz .Fakat ben her mikroserviste
//Farklı bir algoritma kullanmaya çalışıyorum farklı senaryolar düzenliyorum ...
