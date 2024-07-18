package models

type Code struct {
	User   User `gorm:"foreignKey:UserID"`
	UserID int  `json:"userid"`
	Code   int  `json:"code"`
}

type DeleteCode struct {
	User   User `gorm:"foreignKey:UserID"`
	UserID int  `json:"userid"`
	Code   int  `json:"code"`
}

type InputCode struct {
	SendingCode   int `json:"sendingcode"`
	UserInputCode int `json:"userinputcode"`
}
