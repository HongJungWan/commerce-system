package request

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
