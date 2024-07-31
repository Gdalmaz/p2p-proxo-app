package models

type VerifyUpdate struct {
	User   User `gorm:"foreignKey:UserID"`
	UserID int  `json:"userid"`
	Code   int  `json:"code"`
}

type VerifyDelete struct {
	User   User `gorm:"foreignKey:UserID"`
	UserID int  `json:"userid"`
	Code   int  `json:"code"`
}

type Code struct {
	SystemCode int `json:"systemcode"`
	InputCode  int `json:"inputcode"`
}
