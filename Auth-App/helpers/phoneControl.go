package helpers

import (
	"auth/Auth-App/database"
	"auth/Auth-App/models"
)

func PhoneControl(phone string) (bool, error) {
	user := new(models.User)

	err := database.DB.Db.Where("phone=?", phone).First(&user).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
