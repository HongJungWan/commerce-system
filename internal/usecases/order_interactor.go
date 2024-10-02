package usecases

import (
	"errors"
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

func (oi *OrderInteractor) CreateOrder(req *request.CreateOrderRequest, memberNumber string) (*response.CreateOrderResponse, error) {
	order, err := req.ToEntity(memberNumber)
	if err != nil {
		return nil, err
	}

	member, err := oi.MemberRepository.GetByMemberNumber(memberNumber)
	if err != nil || member == nil {
		return nil, errors.New("유효하지 않은 회원 번호입니다.")
	}

	product, err := oi.ProductRepository.GetByProductNumber(order.ProductNumber)
	if err != nil || product == nil {
		return nil, errors.New("유효하지 않은 상품 번호입니다.")
	}

	if product.StockQuantity < order.Quantity {
		return nil, errors.New("재고 수량이 부족합니다.")
	}
	product.StockQuantity -= order.Quantity
	if err := oi.ProductRepository.Update(product); err != nil {
		return nil, err
	}

	order.Price = product.Price
	order.TotalAmount = product.Price * int64(order.Quantity)
	order.IsCanceled = false

	if err := oi.OrderRepository.Create(order); err != nil {
		return nil, err
	}

	orderResponse := response.NewOrderResponse(order)

	return &response.CreateOrderResponse{
		Message: "주문이 등록되었습니다.",
		Order:   *orderResponse,
	}, nil
}

func (oi *OrderInteractor) GetMyOrders(memberNumber string) ([]response.OrderResponse, error) {
	orders, err := oi.OrderRepository.GetByMemberNumber(memberNumber)
	if err != nil {
		return nil, err
	}

	var orderResponses []response.OrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, *response.NewOrderResponse(order))
	}

	return orderResponses, nil
}

func (oi *OrderInteractor) CancelOrder(orderId int, memberNumber string) error {
	order, err := oi.OrderRepository.GetById(orderId)
	if err != nil {
		return err
	}

	if order.MemberNumber != memberNumber {
		return errors.New("해당 주문에 대한 권한이 없습니다.")
	}

	if err := order.Cancel(); err != nil {
		return err
	}

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
