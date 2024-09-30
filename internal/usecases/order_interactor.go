package usecases

import (
	"errors"
	"time"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/domain/repository"
)

type OrderInteractor struct {
	OrderRepository   repository.OrderRepository
	MemberRepository  repository.MemberRepository
	ProductRepository repository.ProductRepository
}

func NewOrderInteractor(or repository.OrderRepository, mr repository.MemberRepository, pr repository.ProductRepository) *OrderInteractor {
	return &OrderInteractor{
		OrderRepository:   or,
		MemberRepository:  mr,
		ProductRepository: pr,
	}
}

func (oi *OrderInteractor) CreateOrder(order *domain.Order) error {
	if err := order.Validate(); err != nil {
		return err
	}

	// 회원 확인
	member, err := oi.MemberRepository.GetByMemberNumber(order.MemberNumber)
	if err != nil || member == nil {
		return errors.New("유효하지 않은 회원 번호입니다.")
	}

	// 상품 확인
	product, err := oi.ProductRepository.GetByProductNumber(order.ProductNumber)
	if err != nil || product == nil {
		return errors.New("유효하지 않은 상품 번호입니다.")
	}

	// 재고 확인 및 감소
	if product.StockQuantity < order.Quantity {
		return errors.New("재고 수량이 부족합니다.")
	}
	product.StockQuantity -= order.Quantity
	if err := oi.ProductRepository.Update(product); err != nil {
		return err
	}

	// 주문 가격 설정
	order.Price = product.Price
	order.TotalAmount = product.Price * int64(order.Quantity)
	order.OrderDate = time.Now()
	order.IsCanceled = true

	return oi.OrderRepository.Create(order)
}

func (oi *OrderInteractor) GetMyOrders(memberNumber string) ([]*domain.Order, error) {
	return oi.OrderRepository.GetByMemberNumber(memberNumber)
}

func (oi *OrderInteractor) CancelOrder(orderNumber, memberNumber string) error {
	order, err := oi.OrderRepository.GetByOrderNumber(orderNumber)
	if err != nil {
		return err
	}

	if order.MemberNumber != memberNumber {
		return errors.New("해당 주문에 대한 권한이 없습니다.")
	}

	if err := order.Cancel(); err != nil {
		return err
	}

	// 재고 복구
	product, err := oi.ProductRepository.GetByProductNumber(order.ProductNumber)
	if err != nil {
		return err
	}
	product.StockQuantity += order.Quantity
	if err := oi.ProductRepository.Update(product); err != nil {
		return err
	}

	return oi.OrderRepository.Update(order)
}

func (oi *OrderInteractor) GetMonthlyStats(month string) (int64, int64, error) {
	return oi.OrderRepository.GetMonthlyStats(month)
}
