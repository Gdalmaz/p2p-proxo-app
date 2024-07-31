package models

type Session struct {
	UserID   int    `gorm:"primaryKey;autoIncrement`
	Token    string `json:"token"`
	IsActive bool   `gorm:"default:true"`
}
