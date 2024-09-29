package domain

import (
	"errors"

	"gorm.io/gorm"
)

type Product struct {
	ProductNumber string `gorm:"primaryKey;column:product_number" json:"product_number"`
	Name          string `gorm:"not null" json:"name"`
	Category      string `json:"category"`
	Price         int64  `gorm:"not null" json:"price"`
	StockQuantity int    `gorm:"not null" json:"stock_quantity"`
}

// 상품 생성 시 유효성 검사
func (p *Product) Validate() error {
	if p.ProductNumber == "" || p.Name == "" || p.Price <= 0 || p.StockQuantity < 0 {
		return errors.New("필수 필드가 누락되었거나 잘못된 값입니다.")
	}
	return nil
}

// 재고 수량 업데이트
func (p *Product) UpdateStock(quantity int) error {
	if quantity < 0 {
		return errors.New("재고 수량은 음수일 수 없습니다.")
	}
	p.StockQuantity = quantity
	return nil
}

// 상품 삭제 가능 여부 확인 (주문 이력이 있는지 체크)
func (p *Product) CanBeDeleted(db *gorm.DB) (bool, error) {
	var count int64
	if err := db.Model(&Order{}).Where("product_number = ?", p.ProductNumber).Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}
