package repository

import "github.com/HongJungWan/commerce-system/internal/domain"

type ProductRepository interface {
	Create(product *domain.Product) error
	GetAll(filter map[string]interface{}) ([]*domain.Product, error)
	GetById(id string) (*domain.Product, error)
	Update(product *domain.Product) error
	Delete(productNumber string) error
}
