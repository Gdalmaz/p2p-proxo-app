package models

type Product struct {
	ID            int    `json:"id"`
	User          User   `gorm:"foreignKey:UserID"`
	UserID        int    `json:"userid"`
	Club          Club   `gorm:"foreignKey:ClubID"`
	ClubID        int    `json:"clubid"`
	ProductName   string `json:"productname"`
	ProductPrice  int    `json:"productprice"`
	ImageUrl      string `json:"imageurl"`
	Image         string `json:"image"`
	ProductAmount int    `json:"productamount"`
}
