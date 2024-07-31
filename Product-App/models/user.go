package models

type User struct {
	ID       int    `json:"id"`
	FullName string `json:"fullname"`
	Password string `json:"password"`
	Mail     string `json:"mail`
	Phone    string `json:string`
	IsActive bool   `gorm:"default:false"`
}

