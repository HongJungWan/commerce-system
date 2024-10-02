package domain

import (
	"errors"

	"gorm.io/gorm"
)

type Product struct {
	ID            int    `gorm:"primaryKey;autoIncrement" json:"id"`                           // 기본 키
	ProductNumber string `gorm:"unique;not null" json:"product_number"`                        // 상품번호
	ProductName   string `gorm:"not null;index:idx_category_product_name" json:"product_name"` // 상품명
	Category      string `gorm:"index:idx_category_product_name" json:"category"`              // 카테고리
	Price         int64  `gorm:"not null" json:"price"`                                        // 가격
	StockQuantity int    `gorm:"not null" json:"stock_quantity"`                               // 재고수량
}

func (p *Product) Validate() error {
	if p.ProductNumber == "" {
		return errors.New("상품번호가 누락되었습니다.")
	}
	if p.ProductName == "" {
		return errors.New("상품명이 누락되었습니다.")
	}
	if p.Price <= 0 {
		return errors.New("가격이 잘못되었습니다.")
	}
	if p.StockQuantity < 0 {
		return errors.New("재고 수량은 음수일 수 없습니다.")
	}
	return nil
}

func (p *Product) UpdateStock(quantity int) error {
	if quantity < 0 {
		return errors.New("재고 수량은 음수일 수 없습니다.")
	}
	p.StockQuantity = quantity
	return nil
}

func (p *Product) CanBeDeleted(db *gorm.DB) (bool, error) {
	var count int64
	if err := db.Model(&Order{}).Where("id = ?", p.ID).Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}
