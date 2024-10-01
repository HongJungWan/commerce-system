package response

import "github.com/HongJungWan/commerce-system/internal/domain"

type ProductResponse struct {
	ID            int    `json:"id"`
	ProductNumber string `json:"product_number"`
	ProductName   string `json:"product_name"`
	Category      string `json:"category"`
	Price         int64  `json:"price"`
	StockQuantity int    `json:"stock_quantity"`
}

type CreateProductResponse struct {
	Message string          `json:"message"`
	Product ProductResponse `json:"product"`
}

func NewProductResponse(product *domain.Product) *ProductResponse {
	return &ProductResponse{
		ID:            product.ID,
		ProductNumber: product.ProductNumber,
		ProductName:   product.ProductName,
		Category:      product.Category,
		Price:         product.Price,
		StockQuantity: product.StockQuantity,
	}
}
