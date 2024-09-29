package domain

import (
	"errors"
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
	Status        string     `gorm:"not null;default:'ordered'" json:"status"`
	CanceledAt    *time.Time `json:"canceled_at,omitempty"`
}

// 주문 생성 시 유효성 검사
func (o *Order) Validate() error {
	if o.OrderNumber == "" || o.MemberNumber == "" || o.ProductNumber == "" || o.Quantity <= 0 || o.Price <= 0 {
		return errors.New("필수 필드가 누락되었거나 잘못된 값입니다.")
	}
	return nil
}

// 주문 취소 처리
func (o *Order) Cancel() error {
	if o.Status != "ordered" {
		return errors.New("이미 취소된 주문입니다.")
	}
	o.Status = "canceled"
	now := time.Now()
	o.CanceledAt = &now
	return nil
}
