package domain

import (
	"errors"
	"time"
)

type Order struct {
	ID            int        `gorm:"primaryKey;autoIncrement" json:"id"`  // 기본 키
	OrderNumber   string     `gorm:"unique;not null" json:"order_number"` // 주문번호
	OrderDate     time.Time  `gorm:"not null" json:"order_date"`          // 주문일
	MemberNumber  string     `gorm:"not null" json:"member_number"`       // 회원번호
	ProductNumber int        `gorm:"not null" json:"product_number"`      // 상품번호
	Price         int64      `gorm:"not null" json:"price"`               // 가격
	Quantity      int        `gorm:"not null" json:"quantity"`            // 수량
	TotalAmount   int64      `gorm:"not null" json:"total_amount"`        // 금액
	IsCanceled    bool       `gorm:"default:false" json:"is_canceled"`    // 취소여부
	CanceledAt    *time.Time `json:"canceled_at,omitempty"`               // 취소일
}

func (o *Order) Validate() error {
	if o.OrderNumber == "" {
		return errors.New("주문번호가 누락되었습니다.")
	}
	if o.OrderDate.IsZero() {
		return errors.New("주문일이 누락되었습니다.")
	}
	if o.MemberNumber == "" {
		return errors.New("회원번호가 누락되었습니다.")
	}
	if o.ProductNumber <= 0 {
		return errors.New("상품번호가 누락되었습니다.")
	}
	if o.Price <= 0 {
		return errors.New("가격이 잘못되었습니다.")
	}
	if o.Quantity <= 0 {
		return errors.New("수량이 잘못되었습니다.")
	}
	if o.TotalAmount <= 0 {
		return errors.New("금액이 잘못되었습니다.")
	}
	return nil
}

func (o *Order) Cancel() error {
	if o.IsCanceled {
		return errors.New("이미 취소된 주문입니다.")
	}
	o.IsCanceled = true
	now := time.Now()
	o.CanceledAt = &now
	return nil
}
