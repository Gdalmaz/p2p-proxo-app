package models

type ForgotPassword struct {
	ID     int    `json:"id"`
	User   User   `gorm:"foreignKey:UserID"`
	UserID int    `json:"userid"`
	Mail   string `json:"mail"`
	Code   int    `json:"code"`
	Token  string `json:"token"`
}

type UpdateForgottenPassword struct {
	Password1 string `json:"password1"`
	Password2 string `json:"password2"`
	Code      int    `json:"code"`
}
