package controller

import (
	"net/http"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/usecases"
	"github.com/gin-gonic/gin"
)

type OrderController struct {
	orderInteractor *usecases.OrderInteractor
}

func NewOrderController(oi *usecases.OrderInteractor) *OrderController {
	return &OrderController{orderInteractor: oi}
}

// 주문 등록
func (oc *OrderController) CreateOrder(c *gin.Context) {
	var order domain.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다."})
		return
	}

	memberNumber := c.GetString("member_number")
	order.MemberNumber = memberNumber

	if err := oc.orderInteractor.CreateOrder(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "주문이 등록되었습니다.", "order": order})
}

// 내 주문 조회
func (oc *OrderController) GetMyOrders(c *gin.Context) {
	memberNumber := c.GetString("member_number")
	orders, err := oc.orderInteractor.GetMyOrders(memberNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "주문 내역을 가져올 수 없습니다."})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// 내 주문 취소
func (oc *OrderController) CancelOrder(c *gin.Context) {
	orderNumber := c.Param("order_number")
	memberNumber := c.GetString("member_number")

	if err := oc.orderInteractor.CancelOrder(orderNumber, memberNumber); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "주문이 취소되었습니다."})
}

// 주문 현황 조회 (월별 매출액, 취소액)
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
