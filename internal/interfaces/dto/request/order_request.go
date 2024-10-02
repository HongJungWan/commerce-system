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
	ProductNumber string `json:"product_number" example:"Product0fe0dfb2-0a9e-4e47-b670-5d1a761e62b5"`
	Quantity      int    `json:"quantity" example:"2"`
	Price         int64  `json:"price" example:"1000"`
}

type CancelOrderRequest struct {
	OrderNumber string `json:"order_number" example:"Order1234567890"`
}

func (req *CreateOrderRequest) CreateToEntity(memberNumber string) (*domain.Order, error) {
	order := &domain.Order{
		OrderNumber:   ORDER + uuid.New().String(),
		OrderDate:     time.Now(),
		MemberNumber:  memberNumber,
		ProductNumber: req.ProductNumber,
		Price:         req.Price,
		Quantity:      req.Quantity,
		TotalAmount:   req.Price * int64(req.Quantity),
	}

	if err := order.Validate(); err != nil {
		return nil, err
	}

	return order, nil
}
