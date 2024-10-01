package request

type CreateOrderRequest struct {
	OrderNumber   string `json:"order_number"`
	ProductNumber string `json:"product_number"`
	Quantity      int    `json:"quantity"`
}

type CancelOrderRequest struct {
	OrderNumber string `json:"order_number"`
}
