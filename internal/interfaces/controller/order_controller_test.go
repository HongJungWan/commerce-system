package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/infrastructure/repository"
	"github.com/HongJungWan/commerce-system/internal/interfaces/controller"
	"github.com/HongJungWan/commerce-system/internal/usecases"
	"github.com/HongJungWan/commerce-system/test/fixtures"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestOrderController_CreateOrder_Failure_InvalidRequest(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	orderRepo := repository.NewOrderRepository(db)
	memberRepo := repository.NewMemberRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderInteractor := usecases.NewOrderInteractor(orderRepo, memberRepo, productRepo)
	orderController := controller.NewOrderController(orderInteractor)

	router := gin.Default()
	router.POST("/orders", func(c *gin.Context) {
		c.Set("member_number", "M12345")
		orderController.CreateOrder(c)
	})

	invalidOrder := map[string]interface{}{
		"quantity": "invalid", // 잘못된 타입의 값
	}
	requestBody, _ := json.Marshal(invalidOrder)
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestOrderController_GetMyOrders_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	orderRepo := repository.NewOrderRepository(db)
	memberRepo := repository.NewMemberRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderInteractor := usecases.NewOrderInteractor(orderRepo, memberRepo, productRepo)
	orderController := controller.NewOrderController(orderInteractor)

	order1 := &domain.Order{
		OrderNumber:   "O12345",
		MemberNumber:  "M12345",
		ProductNumber: 12345,
	}
	order2 := &domain.Order{
		OrderNumber:   "O12346",
		MemberNumber:  "M12345",
		ProductNumber: 12346,
	}
	_ = orderRepo.Create(order1)
	_ = orderRepo.Create(order2)

	router := gin.Default()
	router.GET("/orders", func(c *gin.Context) {
		c.Set("member_number", "M12345")
		orderController.GetMyOrders(c)
	})

	req, _ := http.NewRequest("GET", "/orders", nil)

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusOK, resp.Code)
	var orders []domain.Order
	err := json.Unmarshal(resp.Body.Bytes(), &orders)
	assert.NoError(t, err)
	assert.Len(t, orders, 2)
}

func TestOrderController_CancelOrder_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	orderRepo := repository.NewOrderRepository(db)
	memberRepo := repository.NewMemberRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderInteractor := usecases.NewOrderInteractor(orderRepo, memberRepo, productRepo)
	orderController := controller.NewOrderController(orderInteractor)

	product := &domain.Product{
		ID:            12345,
		ProductNumber: "P12345",
		StockQuantity: 10,
	}
	order := &domain.Order{
		ID:            12345,
		OrderNumber:   "O12345",
		MemberNumber:  "M12345",
		ProductNumber: 12345,
		Quantity:      2,
		IsCanceled:    false,
	}
	_ = productRepo.Create(product)
	_ = orderRepo.Create(order)

	router := gin.Default()
	router.PUT("/orders/:order_number/cancel", func(c *gin.Context) {
		c.Set("member_number", "M12345")
		orderController.CancelOrder(c)
	})

	req, _ := http.NewRequest("PUT", "/orders/O12345/cancel", nil)

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusOK, resp.Code)
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "주문이 취소되었습니다.", response["message"])
}

func TestOrderController_GetMonthlyStats_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	orderRepo := repository.NewOrderRepository(db)
	memberRepo := repository.NewMemberRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderInteractor := usecases.NewOrderInteractor(orderRepo, memberRepo, productRepo)
	orderController := controller.NewOrderController(orderInteractor)

	order1 := &domain.Order{
		ID:            12345,
		OrderNumber:   "O12345",
		OrderDate:     time.Date(2024, 9, 10, 0, 0, 0, 0, time.UTC),
		MemberNumber:  "M12345",
		ProductNumber: 12345,
		Price:         1000,
		Quantity:      2,
		TotalAmount:   2000,
		IsCanceled:    false,
	}
	order2 := &domain.Order{
		ID:            12346,
		OrderNumber:   "O12346",
		OrderDate:     time.Date(2024, 9, 15, 0, 0, 0, 0, time.UTC),
		MemberNumber:  "M12346",
		ProductNumber: 12346,
		Price:         1500,
		Quantity:      1,
		TotalAmount:   1500,
		IsCanceled:    true,
	}
	_ = orderRepo.Create(order1)
	_ = orderRepo.Create(order2)

	router := gin.Default()
	router.GET("/orders/stats", func(c *gin.Context) {
		c.Set("is_admin", true)
		c.Request.URL.RawQuery = "month=2024-09"
		orderController.GetMonthlyStats(c)
	})

	req, _ := http.NewRequest("GET", "/orders/stats", nil)

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusOK, resp.Code)
	var stats map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &stats)
	assert.NoError(t, err)
	assert.Equal(t, float64(2000), stats["total_sales"])
	assert.Equal(t, float64(1500), stats["total_canceled"])
}

func TestOrderController_GetMonthlyStats_Failure_Unauthorized(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	orderRepo := repository.NewOrderRepository(db)
	memberRepo := repository.NewMemberRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderInteractor := usecases.NewOrderInteractor(orderRepo, memberRepo, productRepo)
	orderController := controller.NewOrderController(orderInteractor)

	router := gin.Default()
	router.GET("/orders/stats", func(c *gin.Context) {
		c.Set("is_admin", false)
		orderController.GetMonthlyStats(c)
	})

	req, _ := http.NewRequest("GET", "/orders/stats", nil)

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusForbidden, resp.Code)
}
