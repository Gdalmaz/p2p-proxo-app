package helpers

import (
	"recipe/database"
	"recipe/models"
)

func CheckVerifyUser(ID int) (bool, error) {
	var user models.User
	err := database.DB.Db.Where("id=?", ID).First(&user).Error
	if err != nil {
		return false, err
	}
	if user.IsActive == false {
		return false, err
	}
	return true, nil

}
