package request

import (
	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/google/uuid"
	"time"
)

const (
	ORDER = "Order"
)

type CreateOrderRequest struct {
	ProductNumber string `json:"product_number"`
	Quantity      int    `json:"quantity"`
}

type CancelOrderRequest struct {
	OrderNumber string `json:"order_number"`
}

func (req *CreateOrderRequest) CreateToEntity(memberNumber string) (*domain.Order, error) {
	order := &domain.Order{
		OrderNumber:   ORDER + uuid.New().String(),
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
