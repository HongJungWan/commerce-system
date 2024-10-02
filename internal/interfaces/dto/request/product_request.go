package request

import (
	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/google/uuid"
)

const (
	PRODUCT = "Product"
)

type CreateProductRequest struct {
	ProductName   string `json:"product_name" example:"pizza"`
	Category      string `json:"category" example:"food"`
	Price         int64  `json:"price" example:"1000"`
	StockQuantity int    `json:"stock_quantity" example:"100"`
}

type UpdateStockRequest struct {
	StockQuantity int `json:"stock_quantity" example:"77"`
}

func (req *CreateProductRequest) CreateToEntity() (*domain.Product, error) {
	product := &domain.Product{
		ProductNumber: PRODUCT + uuid.New().String(),
		ProductName:   req.ProductName,
		Category:      req.Category,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
	}

	if err := product.Validate(); err != nil {
		return nil, err
	}

	return product, nil
}
