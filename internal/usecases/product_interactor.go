package usecases

import (
	"errors"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/domain/repository"
	"github.com/HongJungWan/commerce-system/internal/interfaces/dto/request"
	"github.com/HongJungWan/commerce-system/internal/interfaces/dto/response"
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

func (pi *ProductInteractor) CreateProduct(req *request.CreateProductRequest) (*response.CreateProductResponse, error) {
	product := pi.toEntity(req)
	if err := product.Validate(); err != nil {
		return nil, err
	}
	if err := pi.ProductRepository.Create(product); err != nil {
		return nil, err
	}
	return &response.CreateProductResponse{
		Message: "상품이 등록되었습니다.",
		Product: pi.toDTO(product),
	}, nil
}

func (pi *ProductInteractor) GetProducts(filter map[string]interface{}) ([]*response.ProductResponse, error) {
	products, err := pi.ProductRepository.GetAll(filter)
	if err != nil {
		return nil, err
	}

	var productResponses []*response.ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, pi.toDTO(product))
	}

	return productResponses, nil
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

func (pi *ProductInteractor) toEntity(req *request.CreateProductRequest) *domain.Product {
	return &domain.Product{
		ProductNumber: req.ProductNumber,
		ProductName:   req.ProductName,
		Category:      req.Category,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
	}
}

func (pi *ProductInteractor) toDTO(product *domain.Product) *response.ProductResponse {
	return &response.ProductResponse{
		ID:            product.ID,
		ProductNumber: product.ProductNumber,
		ProductName:   product.ProductName,
		Category:      product.Category,
		Price:         product.Price,
		StockQuantity: product.StockQuantity,
	}
}
