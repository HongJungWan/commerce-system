package domain

type Product struct {
	ProductNumber string  `gorm:"primaryKey;column:product_number" json:"product_number"`
	Name          string  `gorm:"not null" json:"name"`
	Category      string  `json:"category"`
	Price         int64   `gorm:"not null" json:"price"`
	StockQuantity int     `gorm:"not null" json:"stock_quantity"`
	Orders        []Order `gorm:"foreignKey:ProductNumber;references:ProductNumber" json:"orders"`
}
