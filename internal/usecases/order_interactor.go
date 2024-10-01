package usecases

import (
	"errors"
	"time"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/domain/repository"
	"github.com/HongJungWan/commerce-system/internal/interfaces/dto/request"
	"github.com/HongJungWan/commerce-system/internal/interfaces/dto/response"
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

// CreateOrder는 주문을 생성하고 응답 DTO를 반환합니다.
func (oi *OrderInteractor) CreateOrder(req *request.CreateOrderRequest, memberNumber string) (*response.CreateOrderResponse, error) {
	// 요청 DTO를 도메인 엔티티로 변환
	order, err := oi.toEntity(req, memberNumber)
	if err != nil {
		return nil, err
	}

	// 회원 확인
	member, err := oi.MemberRepository.GetByMemberNumber(memberNumber)
	if err != nil || member == nil {
		return nil, errors.New("유효하지 않은 회원 번호입니다.")
	}

	// 상품 확인
	product, err := oi.ProductRepository.GetByProductNumber(order.ProductNumber)
	if err != nil || product == nil {
		return nil, errors.New("유효하지 않은 상품 번호입니다.")
	}

	// 재고 확인 및 감소
	if product.StockQuantity < order.Quantity {
		return nil, errors.New("재고 수량이 부족합니다.")
	}
	product.StockQuantity -= order.Quantity
	if err := oi.ProductRepository.Update(product); err != nil {
		return nil, err
	}

	// 주문 가격 설정
	order.Price = product.Price
	order.TotalAmount = product.Price * int64(order.Quantity)
	order.OrderDate = time.Now()
	order.IsCanceled = false

	// 주문 저장
	if err := oi.OrderRepository.Create(order); err != nil {
		return nil, err
	}

	// 도메인 엔티티를 응답 DTO로 변환
	orderResponse := oi.toDTO(order)

	return &response.CreateOrderResponse{
		Message: "주문이 등록되었습니다.",
		Order:   *orderResponse,
	}, nil
}

// GetMyOrders는 회원의 주문 목록을 응답 DTO로 반환합니다.
func (oi *OrderInteractor) GetMyOrders(memberNumber string) ([]response.OrderResponse, error) {
	orders, err := oi.OrderRepository.GetByMemberNumber(memberNumber)
	if err != nil {
		return nil, err
	}

	var orderResponses []response.OrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, *oi.toDTO(order))
	}

	return orderResponses, nil
}

// CancelOrder는 주문을 취소합니다.
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

// GetMonthlyStats는 월별 통계 정보를 반환합니다.
func (oi *OrderInteractor) GetMonthlyStats(month string) (int64, int64, error) {
	return oi.OrderRepository.GetMonthlyStats(month)
}

// 요청 DTO를 도메인 엔티티로 변환
func (oi *OrderInteractor) toEntity(req *request.CreateOrderRequest, memberNumber string) (*domain.Order, error) {
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

// 도메인 엔티티를 응답 DTO로 변환
func (oi *OrderInteractor) toDTO(order *domain.Order) *response.OrderResponse {
	return &response.OrderResponse{
		ID:            order.ID,
		OrderNumber:   order.OrderNumber,
		OrderDate:     order.OrderDate.Format(time.RFC3339),
		MemberNumber:  order.MemberNumber,
		ProductNumber: order.ProductNumber,
		Price:         order.Price,
		Quantity:      order.Quantity,
		TotalAmount:   order.TotalAmount,
		IsCanceled:    order.IsCanceled,
		CanceledAt:    formatTime(order.CanceledAt),
	}
}
