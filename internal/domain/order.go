package domain

import (
	"time"
)

type Order struct {
	OrderNumber   string     `gorm:"primaryKey;column:order_number" json:"order_number"`
	OrderDate     time.Time  `gorm:"not null" json:"order_date"`
	MemberNumber  string     `gorm:"not null" json:"member_number"`
	ProductNumber string     `gorm:"not null" json:"product_number"`
	Price         int64      `gorm:"not null" json:"price"`
	Quantity      int        `gorm:"not null" json:"quantity"`
	TotalPrice    int64      `gorm:"not null" json:"total_price"`
	Canceled      bool       `gorm:"not null;default:false" json:"canceled"`
	CanceledAt    *time.Time `json:"canceled_at"`
}
