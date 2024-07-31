package helpers

import (
	"market/database"
	"market/models"

	"gorm.io/gorm"
)

func ClubControl(userid int) (bool, error) {
	var club models.Club

	err := database.DB.Db.Where("user_id = ?", userid).First(&club).Error
	if err != nil {
		// Eğer hata kayıt bulunamamasından kaynaklanıyorsa
		if gorm.ErrRecordNotFound == err {
			return false, nil
		}
		// Diğer hatalar için
		return false, err
	}
	return true, nil
}