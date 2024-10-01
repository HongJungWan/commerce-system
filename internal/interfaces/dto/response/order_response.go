package response

type OrderResponse struct {
	ID            int    `json:"id"`
	OrderNumber   string `json:"order_number"`
	OrderDate     string `json:"order_date"`
	MemberNumber  string `json:"member_number"`
	ProductNumber string `json:"product_number"`
	Price         int64  `json:"price"`
	Quantity      int    `json:"quantity"`
	TotalAmount   int64  `json:"total_amount"`
	IsCanceled    bool   `json:"is_canceled"`
	CanceledAt    string `json:"canceled_at,omitempty"`
}

type CreateOrderResponse struct {
	Message string        `json:"message"`
	Order   OrderResponse `json:"order"`
}
