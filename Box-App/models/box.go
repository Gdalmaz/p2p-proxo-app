package models

type Box struct {
	ID        int    `json:"id"` //BOX ID OLARAK KULLANACAĞIZ
	User      User   `gorm:"foreignKey:UserID"`
	UserID    int    `json:"userid"`
	Adress    string `json:"adress"`
	Phone     string `json:"phone"`
	TotalCost int    `json:"totalcost" gorm:"default:0"`
}

type AddProduct struct {
	Box           Box     `gorm:"foreignKey:BoxID"`
	BoxID         int     `json:"boxid"`
	User          User    `json:"UserID"`
	UserID        int     `json:"userid"`
	Product       Product `gorm:"foreignKey:ProductID"`
	ProductID     int     `json:"productid"`
	Queue         int     `gorm:"default:1"`
	ProductAmount int     `json:"productamount"`
}

// ADD PRODUCT KISMINDA Kİ ALGORİTMAMIZ TAM OLARAK ŞÖYLE OLUCAK :

//KULLANICI SEPETE ÜRÜN EKLEYECEK BOX İD KENDİLİĞİNDEN OLUŞUCAK KULLANICI İLK ALIŞVERİŞİNDE KULLANICIDAN BİLGİ İSTEYECEĞİZ KONUMUNU ADRESİNİ TEKRARDAN TELEFONUNU İSTEYECEĞİZ
//BİLGİLER DOĞRULANDIKTAN SONRA ALIŞVERİŞ SEPETİNİ OLUŞTURACAĞIZ VE BUNU DA REDİSE ATACAĞIZ BOX SADECE BİR ARACI ASIL İŞLEM YAPACAĞIMIZ KISIM ADD PRODUCT KISMI OLUCAK .
