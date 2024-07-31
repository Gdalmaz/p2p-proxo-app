package models

import "time"

type Club struct {
	ID           int       `json:"id"`
	User         User      `gorm:"foreignKey:UserID"`
	UserID       int       `json:"userid"`
	ClubName     string    `json:"clubname"`
	Explanation  string    `json:"explanation"`
	UserFullName string    `json:"userfullname"`
	ImageUrl     string    `json:"imageurl"`
	Image        string    `json:"image"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
