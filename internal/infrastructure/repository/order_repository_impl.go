package repository

import (
	"time"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"gorm.io/gorm"
)

type OrderRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepositoryImpl {
	return &OrderRepositoryImpl{db: db}
}

func (r *OrderRepositoryImpl) Create(order *domain.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepositoryImpl) GetByOrderNumber(orderNumber string) (*domain.Order, error) {
	var order domain.Order
	if err := r.db.First(&order, "order_number = ?", orderNumber).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepositoryImpl) GetByMemberNumber(memberNumber string) ([]*domain.Order, error) {
	var orders []*domain.Order
	if err := r.db.Where("member_number = ?", memberNumber).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepositoryImpl) Update(order *domain.Order) error {
	return r.db.Save(order).Error
}

func (r *OrderRepositoryImpl) GetMonthlyStats(month string) (int64, int64, error) {
	startDate, err := time.Parse("2006-01", month)
	if err != nil {
		return 0, 0, err
	}
	endDate := startDate.AddDate(0, 1, 0)

	var totalSales int64
	var totalCanceled int64

	// 매출액 계산 (취소되지 않은 주문)
	if err := r.db.Model(&domain.Order{}).
		Where("order_date >= ? AND order_date < ? AND status = ?", startDate, endDate, "ordered").
		Select("SUM(total_price)").Scan(&totalSales).Error; err != nil {
		return 0, 0, err
	}

	// 취소액 계산 (취소된 주문)
	if err := r.db.Model(&domain.Order{}).
		Where("order_date >= ? AND order_date < ? AND status = ?", startDate, endDate, "canceled").
		Select("SUM(total_price)").Scan(&totalCanceled).Error; err != nil {
		return 0, 0, err
	}

	return totalSales, totalCanceled, nil
}
