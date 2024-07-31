package models

type Box struct {
	ID       int          `json:"id"`
	User     User         `gorm:"foreignKey:UserID"`
	UserID   int          `json:"userid"`
	Products []AddProduct `gorm:"foreignKey:BoxID"`
}

type AddProduct struct {
	ID            int     `json:"id"`
	BoxID         int     `json:"boxid"`
	Queue         int     `gorm:"default:1"`
	Product       Product `gorm:"foreignKey:ProductID"`
	ProductID     int     `json:"productid"`
	ProductAmount int     `gorm:"default:1"`
}
