package controllers

import (
	"auth/config"
	"auth/database"
	"auth/helpers"
	"auth/middleware"
	"auth/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SignUp(c *fiber.Ctx) error {
	user := new(models.User)
	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "S-U-1"})
	}
	mailControl, _ := helpers.MailControl(user.Mail)
	if mailControl == true {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "S-U-2"})
	}
	phoneControl, _ := helpers.PhoneControl(user.Mail)
	if phoneControl == true {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "S-U-3"})
	}
	user.Password = helpers.HashPass(user.Password)

	err = database.DB.Db.Create(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "S-U-4"})
	}
	var codeGenerator models.Code
	code := helpers.GenerateRandomNumber()
	codeGenerator.Code = code
	codeGenerator.UserID = user.ID
	codeStr := strconv.Itoa(code)
	err = config.RabbitMqPublish([]byte(codeStr), user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "S-U-5", "data": err.Error()})
	}
	err = config.RabbitMqConsume(user.Mail, user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "S-U-6"})
	}
	err = database.DB.Db.Create(&codeGenerator).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "S-U-7"})
	}

	var verifySession models.VerifySession
	token, err := middleware.CreateToken(user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "S-U-8"})
	}
	verifySession.UserID = user.ID
	verifySession.Token = token
	err = database.DB.Db.Create(&verifySession).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "S-U-9"})
	}
	return c.Status(202).JSON(fiber.Map{"status": "Success", "message": "Success"})
}

func LastStepSignUp(c *fiber.Ctx) error {
	verifysession, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "L-S-S-U-1"})
	}

	var inputCode models.InputCode
	err := c.BodyParser(&inputCode)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "L-S-S-U-2"})
	}

	var systemCode models.Code
	err = database.DB.Db.Where("user_id=?", verifysession.ID).First(&systemCode).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "L-S-S-U-3"})
	}
	inputCode.SendingCode = systemCode.Code
	if inputCode.UserInputCode != inputCode.SendingCode {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "L-S-S-U-4"})
	}

	verifysession.IsActive = true
	var searchtabletoken models.VerifySession
	err = database.DB.Db.Where("id=?", verifysession.ID).Save(&verifysession).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "L-S-S-U-5"})
	}
	err = database.DB.Db.Where("user_id=?", verifysession.ID).Delete(&searchtabletoken).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "L-S-S-U-6"})
	}
	err = database.DB.Db.Where("user_id=?", verifysession.ID).Delete(&systemCode).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "L-S-S-U-7"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success", "data": verifysession})

}

func LogIn(c *fiber.Ctx) error {
	var user models.User
	var login models.LogIn
	err := c.BodyParser(&login)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "L-I-1"})
	}
	login.Password = helpers.HashPass(login.Password)
	err = database.DB.Db.Where("mail=? and password=?", login.Mail, login.Password).First(&user).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "L-I-2"})
	}

	var session models.Session
	token, err := middleware.CreateToken(user.Mail)
	session.Token = token
	session.UserID = user.ID
	err = database.DB.Db.Create(&session).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "L-I-3"})
	}

	//KULLANICIYI BİLGİLENDİRME
	text := "login on your profile . LogIn Time : \t" + time.Now().Format("2006-01-02 15:04:05")
	err = config.RabbitMqPublish([]byte(text), user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "L-I-4"})
	}
	err = config.RabbitMqConsume(user.Mail, user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "L-I-5"})
	}

	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success"})
}

func UpdateAccount(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-A-1"})
	}

	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-A-2"})
	}

	err = database.DB.Db.Where("id=?", user.ID).Updates(&user).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "U-A-3"})
	}
	text := "Your profile information updated"
	err = config.RabbitMqPublish([]byte(text), user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-A-4"})
	}
	err = config.RabbitMqConsume(user.Mail, user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-A-5"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "your profile updated successfully", "data": user})
}

func UpdatePassword(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "U-P-1"})
	}
	var updateinfo models.UpdatePassword
	err := c.BodyParser(&updateinfo)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-P-2"})
	}
	user.Password = helpers.HashPass(user.Password)
	updateinfo.OldPassword = user.Password
	if updateinfo.OldPassword != user.Password {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "U-P-3"})
	}
	if updateinfo.NewPassword1 != updateinfo.NewPassword2 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "U-P-4"})
	}

	user.Password = updateinfo.NewPassword1
	sendmailinuserpassword := user.Password
	user.Password = helpers.HashPass(user.Password)
	err = database.DB.Db.Where("id=?", user.ID).Updates(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-P-5"})
	}

	text := "Your password updated . Your new password is : \n\t\t" + sendmailinuserpassword

	err = config.RabbitMqPublish([]byte(text), user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-P-6"})
	}
	err = config.RabbitMqConsume(user.Mail, user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-P-7"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success"})
}

func LogOut(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "L-O-1"})
	}
	var session models.Session
	err := database.DB.Db.Where("user_id=?", user.ID).Delete(&session).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "L-O-2"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success"})
}
func DeleteAccountSendMail(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "D-A-S-M-1"})
	}
	// kullanıcının hesabına mail yollama işlemi güvenlik zaafiyeti
	var sendingCode models.DeleteCode
	code := helpers.GenerateRandomNumber()
	sendingCode.Code = code
	codeStr := strconv.Itoa(code)
	sendingCode.UserID = user.ID
	err := config.RabbitMqPublish([]byte(codeStr), user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "D-A-S-M-2"})
	}
	err = config.RabbitMqConsume(user.Mail, user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "D-A-S-M-3"})
	}
	err = database.DB.Db.Create(&sendingCode).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "D-A-S-M-4"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success"})
}

func DeleteAccountVerifyMail(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "D-A-V-M-1"})
	}
	var sendingCode models.DeleteCode
	err := database.DB.Db.Where("user_id=?", user.ID).Find(&sendingCode).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "D-A-V-M-2"})
	}
	var inputCode models.InputCode
	err = c.BodyParser(&inputCode)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "D-A-V-M-3"})
	}
	inputCode.SendingCode = sendingCode.Code
	if inputCode.UserInputCode != inputCode.SendingCode {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "D-A-V-M-4"})
	}

	text := "Your profile deleted successfully . Thank you for your services " + time.Now().Format("2006-01-02 15:04:05")
	mail := user.Mail

	var session models.Session
	err = database.DB.Db.Where("user_id=?", user.ID).Delete(&session).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "D-A-V-M-5"})
	}
	err = database.DB.Db.Where("user_id=?", user.ID).Delete(&sendingCode).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "D-A-V-M-6"})
	}

	err = database.DB.Db.Where("id=?", user.ID).Delete(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "D-A-V-M-7"})
	}
	err = config.RabbitMqPublish([]byte(text), mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "D-A-V-M-8"})
	}
	err = config.RabbitMqConsume(mail, mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "D-A-V-M-9"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "mesage": "Success"})
}

func ForgotPasswordSendCode(c *fiber.Ctx) error {
	var user models.User
	var forgotpassword models.ForgotPassword
	err := c.BodyParser(&forgotpassword)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "F-P-S-1"})
	}
	err = database.DB.Db.Where("mail=?", forgotpassword.Mail).Find(&user).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "F-P-S-2"})
	}
	code := helpers.GenerateRandomNumber()
	forgotpassword.UserID = user.ID
	forgotpassword.Code = code
	token, err := middleware.CreateToken(forgotpassword.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "F-P-S-3"})
	}
	forgotpassword.Token = token
	codeStr := strconv.Itoa(code)
	err = config.RabbitMqPublish([]byte(codeStr), forgotpassword.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "F-P-S-4"})
	}
	err = config.RabbitMqConsume(forgotpassword.Mail, forgotpassword.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "F-P-S-5"})
	}
	err = database.DB.Db.Create(&forgotpassword).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "F-P-S-6"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success"})
}

func ForgotPasswordVerifyAndReset(c *fiber.Ctx) error {
	// Forgot password bilgilerini al
	forgotpassword, ok := c.Locals("user").(models.ForgotPassword)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "F-P-V-A-R-1"})
	}

	// Body'den updatepassword bilgilerini al
	var updatepassword models.UpdateForgottenPassword
	err := c.BodyParser(&updatepassword)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "F-P-V-A-R-2"})
	}

	// Kod doğrulaması yap
	if forgotpassword.Code != updatepassword.Code {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "F-P-V-A-R-3"})
	}

	// Şifre doğrulaması yap
	if updatepassword.Password1 != updatepassword.Password2 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "F-P-V-A-R-4"})
	}

	// Kullanıcıyı güncelle
	var user models.User
	err = database.DB.Db.Where("id=?", forgotpassword.UserID).First(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "F-P-V-A-R-5"})
	}
	updatepassword.Password1 = helpers.HashPass(updatepassword.Password1)
	user.Password = updatepassword.Password1
	err = database.DB.Db.Updates(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "F-P-V-A-R-6"})
	}

	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success"})
}
