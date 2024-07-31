package helpers

import (
	"box/database"
	"box/models"
	"errors"

	"gorm.io/gorm"
)

func BoxCreateControl(id int) (bool, error) {
	var box models.Box

	err := database.DB.Db.Where("user_id = ?", id).First(&box).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil // Box not found for the given user_id
		}
		return false, err // An error occurred during the query
	}
	return true, nil // Box found
}
