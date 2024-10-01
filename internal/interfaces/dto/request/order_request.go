package request

import (
	"github.com/HongJungWan/commerce-system/internal/domain"
	"time"
)

type CreateOrderRequest struct {
	OrderNumber   string `json:"order_number"`
	ProductNumber string `json:"product_number"`
	Quantity      int    `json:"quantity"`
}

type CancelOrderRequest struct {
	OrderNumber string `json:"order_number"`
}

func (req *CreateOrderRequest) ToEntity(memberNumber string) (*domain.Order, error) {
	order := &domain.Order{
		OrderNumber:   req.OrderNumber,
		OrderDate:     time.Now(),
		MemberNumber:  memberNumber,
		ProductNumber: req.ProductNumber,
		Quantity:      req.Quantity,
	}

	if err := order.Validate(); err != nil {
		return nil, err
	}

	return order, nil
}
