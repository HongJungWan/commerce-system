package request

import (
	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/google/uuid"
)

const (
	PRODUCT = "Product"
)

type CreateProductRequest struct {
	ProductName   string `json:"product_name"`
	Category      string `json:"category"`
	Price         int64  `json:"price"`
	StockQuantity int    `json:"stock_quantity"`
}

type UpdateStockRequest struct {
	StockQuantity int `json:"stock_quantity"`
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
