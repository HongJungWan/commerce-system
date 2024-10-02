package controller

import (
	"net/http"
	"strconv"

	"github.com/HongJungWan/commerce-system/internal/interfaces/dto/request"
	"github.com/HongJungWan/commerce-system/internal/usecases"
	"github.com/gin-gonic/gin"
)

type OrderController struct {
	orderInteractor *usecases.OrderInteractor
}

func NewOrderController(oi *usecases.OrderInteractor) *OrderController {
	return &OrderController{orderInteractor: oi}
}

// CreateOrder godoc
// @Summary      주문 생성
// @Description  새로운 주문을 생성합니다.
// @Tags         orders
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        orderRequest body request.CreateOrderRequest true "주문 정보"
// @Success      201 {object} response.OrderResponse "주문 생성 성공"
// @Failure      400 {object} map[string]string "잘못된 요청"
// @Failure      500 {object} map[string]string "주문 생성 실패"
// @Router       /orders [post]
func (oc *OrderController) CreateOrder(c *gin.Context) {
	var req request.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다."})
		return
	}

	memberNumber := c.GetString("member_number")

	responseData, err := oc.orderInteractor.CreateOrder(&req, memberNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, responseData)
}

// GetMyOrders godoc
// @Summary      내 주문 조회
// @Description  인증된 사용자의 주문 내역을 조회합니다.
// @Tags         orders
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Success      200 {array} response.OrderResponse "주문 목록"
// @Failure      500 {object} map[string]string "주문 조회 실패"
// @Router       /orders/me [get]
func (oc *OrderController) GetMyOrders(c *gin.Context) {
	memberNumber := c.GetString("member_number")

	orderResponses, err := oc.orderInteractor.GetMyOrders(memberNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "주문 내역을 가져올 수 없습니다."})
		return
	}
	c.JSON(http.StatusOK, orderResponses)
}

// CancelOrder godoc
// @Summary      주문 취소
// @Description  특정 주문을 취소합니다.
// @Tags         orders
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        ID path string true "기본키 (primary key)"
// @Success      200 {object} map[string]string "취소 성공"
// @Failure      500 {object} map[string]string "취소 실패"
// @Router       /orders/:id/cancel [put]
func (oc *OrderController) CancelOrder(c *gin.Context) {
	orderParam := c.Param("id")
	id, err := strconv.Atoi(orderParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 주문 ID입니다."})
		return
	}

	memberNumber := c.GetString("member_number")

	if err := oc.orderInteractor.CancelOrder(id, memberNumber); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "주문이 취소되었습니다."})
}

// GetMonthlyStats godoc
// @Summary      주문 통계 조회
// @Description  특정 월의 주문 통계를 조회합니다. (관리자 전용)
// @Tags         orders
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        month query string true "조회할 월 (YYYY-MM)"
// @Success      200 {object} response.OrderStatsResponse "통계 정보"
// @Failure      400 {object} map[string]string "잘못된 요청"
// @Failure      403 {object} map[string]string "권한 없음"
// @Failure      500 {object} map[string]string "통계 조회 실패"
// @Router       /orders/stats [get]
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
