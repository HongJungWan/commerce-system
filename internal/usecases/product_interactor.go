package usecases

import (
	"errors"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/domain/repository"
	"gorm.io/gorm"
)

type ProductInteractor struct {
	ProductRepository repository.ProductRepository
	DB                *gorm.DB
}

func NewProductInteractor(repo repository.ProductRepository, db *gorm.DB) *ProductInteractor {
	return &ProductInteractor{
		ProductRepository: repo,
		DB:                db,
	}
}

func (pi *ProductInteractor) CreateProduct(product *domain.Product) error {
	if err := product.Validate(); err != nil {
		return err
	}
	return pi.ProductRepository.Create(product)
}

func (pi *ProductInteractor) GetProducts(filter map[string]interface{}) ([]*domain.Product, error) {
	return pi.ProductRepository.GetAll(filter)
}

func (pi *ProductInteractor) UpdateStock(id int, quantity int) error {
	product, err := pi.ProductRepository.GetById(id)
	if err != nil {
		return err
	}
	if err := product.UpdateStock(quantity); err != nil {
		return err
	}
	return pi.ProductRepository.Update(product)
}

func (pi *ProductInteractor) DeleteProduct(id int) error {
	product, err := pi.ProductRepository.GetById(id)
	if err != nil {
		return err
	}
	canBeDeleted, err := product.CanBeDeleted(pi.DB)
	if err != nil {
		return err
	}
	if !canBeDeleted {
		return errors.New("주문된 이력이 있어 삭제할 수 없습니다.")
	}
	return pi.ProductRepository.Delete(id)
}
