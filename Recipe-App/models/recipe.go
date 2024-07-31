package models

import "time"

type Recipe struct {
	ID           int       `json:"id"`
	User         User      `gorm:"foreignKey:UserID"`
	UserID       int       `json:"userid"`
	FoodName     string    `json:"foodname"`
	EatCapacity  int       `json:"eatcapacity"`
	PrepeareTime int       `json:"prepearetime"`
	GuessPrice   int       `json:"guessprice"`
	Materials    string    `json:"materials"`
	Description  string    `json:"description"`
	ImageUrl     string    `json:"imageurl"`
	Image        string    `json:"image"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Popularity struct {
	User        User   `gorm:"foreignKey:UserID"`
	Recipe      Recipe `gorm:"foreignKey:FoodID"`
	ID          int    `json:"id"`
	UserID      int    `json:"userid"`
	FoodID      int    `json:"foodid"`
	ClickNumber int    `gorm:"default:0"`
}
