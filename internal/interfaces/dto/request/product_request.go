package request

import "github.com/HongJungWan/commerce-system/internal/domain"

type CreateProductRequest struct {
	ProductNumber string `json:"product_number"`
	ProductName   string `json:"product_name"`
	Category      string `json:"category"`
	Price         int64  `json:"price"`
	StockQuantity int    `json:"stock_quantity"`
}

type UpdateStockRequest struct {
	StockQuantity int `json:"stock_quantity"`
}

func (req *CreateProductRequest) ToEntity() (*domain.Product, error) {
	product := &domain.Product{
		ProductNumber: req.ProductNumber,
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
