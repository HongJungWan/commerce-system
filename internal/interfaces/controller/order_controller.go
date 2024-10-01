package controller

import (
	"net/http"
	"time"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/interfaces/dto/request"
	"github.com/HongJungWan/commerce-system/internal/interfaces/dto/response"
	"github.com/HongJungWan/commerce-system/internal/usecases"
	"github.com/gin-gonic/gin"
)

type OrderController struct {
	orderInteractor *usecases.OrderInteractor
}

func NewOrderController(oi *usecases.OrderInteractor) *OrderController {
	return &OrderController{orderInteractor: oi}
}

func (oc *OrderController) CreateOrder(c *gin.Context) {
	var req request.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다."})
		return
	}

	memberNumber := c.GetString("member_number")

	order := &domain.Order{
		OrderNumber:   req.OrderNumber,
		OrderDate:     time.Now(),
		MemberNumber:  memberNumber,
		ProductNumber: req.ProductNumber,
		Quantity:      req.Quantity,
	}

	// 가격 및 총 금액은 비즈니스 로직에서 계산됩니다.
	if err := oc.orderInteractor.CreateOrder(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	orderResponse := response.OrderResponse{
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

	c.JSON(http.StatusCreated, response.CreateOrderResponse{
		Message: "주문이 등록되었습니다.",
		Order:   orderResponse,
	})
}

func (oc *OrderController) GetMyOrders(c *gin.Context) {
	memberNumber := c.GetString("member_number")
	orders, err := oc.orderInteractor.GetMyOrders(memberNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "주문 내역을 가져올 수 없습니다."})
		return
	}

	var orderResponses []response.OrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, response.OrderResponse{
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
		})
	}

	c.JSON(http.StatusOK, orderResponses)
}

func (oc *OrderController) CancelOrder(c *gin.Context) {
	orderNumber := c.Param("order_number")
	memberNumber := c.GetString("member_number")

	if err := oc.orderInteractor.CancelOrder(orderNumber, memberNumber); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "주문이 취소되었습니다."})
}

func (oc *OrderController) GetMonthlyStats(c *gin.Context) {
	isAdmin := c.GetBool("is_admin")
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "접근 권한이 없습니다."})
		return
	}

	month := c.Query("month")
	if month == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "month 파라미터가 필요합니다."})
		return
	}

	totalSales, totalCanceled, err := oc.orderInteractor.GetMonthlyStats(month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "통계 정보를 가져올 수 없습니다."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"month":          month,
		"total_sales":    totalSales,
		"total_canceled": totalCanceled,
	})
}
