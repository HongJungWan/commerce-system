package repository

import (
	"github.com/HongJungWan/commerce-system/internal/domain"
)

type OrderRepository interface {
	Create(order *domain.Order) error
	GetByOrderNumber(orderNumber string) (*domain.Order, error)
	GetById(id int) (*domain.Order, error)
	GetByMemberNumber(memberNumber string) ([]*domain.Order, error)
	Update(order *domain.Order) error
	GetMonthlyStats(month string) (int64, int64, error)
}
