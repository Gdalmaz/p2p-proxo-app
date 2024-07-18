package controllers

import (
	"auth/Auth-App/config"
	"auth/Auth-App/database"
	"auth/Auth-App/helpers"
	"auth/Auth-App/middleware"
	"auth/Auth-App/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SignUp(c *fiber.Ctx) error {
	user := new(models.User)
	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error bodyparsing step"})
	}
	mailControl, _ := helpers.MailControl(user.Mail)
	if mailControl == true {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "you already have this account on this mail"})
	}
	phoneControl, _ := helpers.PhoneControl(user.Mail)
	if phoneControl == true {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "you already have this account on this phone number"})
	}
	user.Password = helpers.HashPass(user.Password)

	err = database.DB.Db.Create(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... not creating user table", "data": err.Error()})
	}
	var codeGenerator models.Code
	code := helpers.GenerateRandomNumber()
	codeGenerator.Code = code
	codeGenerator.UserID = user.ID
	codeStr := strconv.Itoa(code)
	err = config.RabbitMqPublish([]byte(codeStr), user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error publishing code step"})
	}
	err = config.RabbitMqConsume(user.Mail, user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error consume code step"})
	}
	err = database.DB.Db.Create(&codeGenerator).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error... not creating code table", "data": err.Error()})
	}

	var verifySession models.VerifySession
	token, err := middleware.CreateToken(user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... not creating token", "data": err.Error()})
	}
	verifySession.UserID = user.ID
	verifySession.Token = token
	err = database.DB.Db.Create(&verifySession).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... not creating token data", "data": err.Error()})
	}
	return c.Status(202).JSON(fiber.Map{"status": "succes", "message": "successfully creating your log"})
}

func LastStepSignUp(c *fiber.Ctx) error {
	verifysession, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "haven't got account ... Please Sign Up"})
	}

	var inputCode models.InputCode
	err := c.BodyParser(&inputCode)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error input code bodyparsing step"})
	}

	var systemCode models.Code
	err = database.DB.Db.Where("user_id=?", verifysession.ID).First(&systemCode).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... not founded code"})
	}
	inputCode.SendingCode = systemCode.Code
	if inputCode.UserInputCode != inputCode.SendingCode {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "you have to input sample where we sending code"})
	}

	verifysession.IsActive = true
	var searchtabletoken models.VerifySession
	err = database.DB.Db.Where("id=?", verifysession.ID).Save(&verifysession).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error... update user database", "data": err.Error()})
	}
	err = database.DB.Db.Where("user_id=?", verifysession.ID).Delete(&searchtabletoken).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... your verify token not deleted", "data": err.Error()})
	}
	err = database.DB.Db.Where("user_id=?", verifysession.ID).Delete(&systemCode).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... your code not deleted on system", "data": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "your profile update to active successfully", "data": verifysession})

}

func LogIn(c *fiber.Ctx) error {
	var user models.User
	var login models.LogIn
	err := c.BodyParser(&login)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... your body parser step"})
	}
	login.Password = helpers.HashPass(login.Password)
	err = database.DB.Db.Where("mail=? and password=?", login.Mail, login.Password).First(&user).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "your password or mail wrong"})
	}

	var session models.Session
	token, err := middleware.CreateToken(user.Mail)
	session.Token = token
	session.UserID = user.ID
	err = database.DB.Db.Create(&session).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... your token not created on table", "data": err.Error()})
	}

	//KULLANICIYI BİLGİLENDİRME
	text := "login on your profile . LogIn Time : \t" + time.Now().Format("2006-01-02 15:04:05")
	err = config.RabbitMqPublish([]byte(text), user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... not publishing mail"})
	}
	err = config.RabbitMqConsume(user.Mail, user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... not consuming mail"})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "your login successfully"})
}

func UpdateAccount(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... not founded token"})
	}

	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...not founded token"})
	}

	err = database.DB.Db.Where("id=?", user.ID).Updates(&user).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error ... update fail", "data": err.Error()})
	}
	text := "Your profile information updated"
	err = config.RabbitMqPublish([]byte(text), user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... not publishing your mail"})
	}
	err = config.RabbitMqConsume(user.Mail, user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... not consuming your mail"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "your profile updated successfully", "data": user})
}

func UpdatePassword(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error ... your token not founded"})
	}
	var updateinfo models.UpdatePassword
	err := c.BodyParser(&updateinfo)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... body parsing step is failed"})
	}
	user.Password = helpers.HashPass(user.Password)
	updateinfo.OldPassword = user.Password
	if updateinfo.OldPassword != user.Password {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error ... your passwords not equals each other"})
	}
	if updateinfo.NewPassword1 != updateinfo.NewPassword2 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error ... your new passwords not equals each other"})
	}

	user.Password = updateinfo.NewPassword1
	sendmailinuserpassword := user.Password
	user.Password = helpers.HashPass(user.Password)
	err = database.DB.Db.Where("id=?", user.ID).Updates(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... your password not updated successfully", "data": err.Error()})
	}

	text := "Your password updated . Your new password is : \n\t\t" + sendmailinuserpassword

	err = config.RabbitMqPublish([]byte(text), user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... not publishing your mail"})
	}
	err = config.RabbitMqConsume(user.Mail, user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... not consuming your mail"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "error", "message": "your password updated successfully"})
}

func LogOut(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "user token not found"})
	}
	var session models.Session
	err := database.DB.Db.Where("user_id=?", user.ID).Delete(&session).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... your token not deleted on table", "data": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "your logout successfully"})
}
func DeleteAccountSendMail(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... your token not founded"})
	}
	// kullanıcının hesabına mail yollama işlemi güvenlik zaafiyeti
	var sendingCode models.DeleteCode
	code := helpers.GenerateRandomNumber()
	sendingCode.Code = code
	codeStr := strconv.Itoa(code)
	sendingCode.UserID = user.ID
	err := config.RabbitMqPublish([]byte(codeStr), user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... error publish mail"})
	}
	err = config.RabbitMqConsume(user.Mail, user.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... error consuming mail"})
	}
	err = database.DB.Db.Create(&sendingCode).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... your code not save "})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "check your mailbox"})
}

func DeleteAccountVerifyMail(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error ... your token not found"})
	}
	var sendingCode models.DeleteCode
	err := database.DB.Db.Where("user_id=?", user.ID).Find(&sendingCode).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "you code not found"})
	}
	var inputCode models.InputCode
	err = c.BodyParser(&inputCode)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... body parser step"})
	}
	inputCode.SendingCode = sendingCode.Code
	if inputCode.UserInputCode != inputCode.SendingCode {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "you have to input sample code which we sending mailbox"})
	}

	text := "Your profile deleted successfully . Thank you for your services " + time.Now().Format("2006-01-02 15:04:05")
	mail := user.Mail

	var session models.Session
	err = database.DB.Db.Where("user_id=?", user.ID).Delete(&session).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... failed to deleting your session"})
	}
	err = database.DB.Db.Where("user_id=?", user.ID).Delete(&sendingCode).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... failed to deleteing your codes"})
	}

	err = database.DB.Db.Where("id=?", user.ID).Delete(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... failed to deleting your data", "data": err.Error()})
	}
	err = config.RabbitMqPublish([]byte(text), mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... failed to publishing your mail"})
	}
	err = config.RabbitMqConsume(mail, mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... failed to consuming your mail"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "mesage": "successfull your delete data"})
}

func ForgotPasswordSendCode(c *fiber.Ctx) error {
	var user models.User
	var forgotpassword models.ForgotPassword
	err := c.BodyParser(&forgotpassword)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "ERROR CODE F-P-S-1"})
	}
	err = database.DB.Db.Where("mail=?", forgotpassword.Mail).Find(&user).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error ... we not found mail for sign users"})
	}
	code := helpers.GenerateRandomNumber()
	forgotpassword.UserID = user.ID
	forgotpassword.Code = code
	token, err := middleware.CreateToken(forgotpassword.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... unsuccessfull creat token"})
	}
	forgotpassword.Token = token
	codeStr := strconv.Itoa(code)
	err = config.RabbitMqPublish([]byte(codeStr), forgotpassword.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... failed to publish your mail"})
	}
	err = config.RabbitMqConsume(forgotpassword.Mail, forgotpassword.Mail)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... failed to consume your mail"})
	}
	err = database.DB.Db.Create(&forgotpassword).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... failed to create steps"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "check your mailbox please"})
}

func ForgotPasswordVerifyAndReset(c *fiber.Ctx) error {
	// Forgot password bilgilerini al
	forgotpassword, ok := c.Locals("user").(models.ForgotPassword)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... not found token"})
	}

	// Body'den updatepassword bilgilerini al
	var updatepassword models.UpdateForgottenPassword
	err := c.BodyParser(&updatepassword)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... error update password body parsing step", "data": err.Error()})
	}

	// Kod doğrulaması yap
	if forgotpassword.Code != updatepassword.Code {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "you should input sample code where we send your mailbox"})
	}

	// Şifre doğrulaması yap
	if updatepassword.Password1 != updatepassword.Password2 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "you should input sample password"})
	}

	// Kullanıcıyı güncelle
	var user models.User
	err = database.DB.Db.Where("id=?", forgotpassword.UserID).First(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... user not found", "data": err.Error()})
	}
	updatepassword.Password1 = helpers.HashPass(updatepassword.Password1)
	user.Password = updatepassword.Password1  // Yeni şifreyi ayarla
	err = database.DB.Db.Updates(&user).Error // Kullanıcıyı güncelle
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error ... failed to update password", "data": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "your password update successfully"})
}
