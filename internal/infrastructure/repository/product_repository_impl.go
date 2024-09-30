package repository

import (
	"github.com/HongJungWan/commerce-system/internal/domain"
	"gorm.io/gorm"
)

type ProductRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepositoryImpl {
	return &ProductRepositoryImpl{db: db}
}

func (r *ProductRepositoryImpl) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepositoryImpl) GetAll(filter map[string]interface{}) ([]*domain.Product, error) {
	var products []*domain.Product
	query := r.db.Model(&domain.Product{})

	if category, ok := filter["category"]; ok {
		query = query.Where("category = ?", category)
	}
	if productName, ok := filter["product_name"]; ok {
		query = query.Where("product_name LIKE ?", "%"+productName.(string)+"%")
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepositoryImpl) GetById(id string) (*domain.Product, error) {
	var product domain.Product
	if err := r.db.First(&product, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepositoryImpl) Update(product *domain.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepositoryImpl) Delete(productNumber string) error {
	return r.db.Delete(&domain.Product{}, "product_number = ?", productNumber).Error
}
