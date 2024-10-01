package usecases_test

import (
	"testing"
	"time"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/infrastructure/repository"
	"github.com/HongJungWan/commerce-system/internal/interfaces/dto/request"
	"github.com/HongJungWan/commerce-system/internal/usecases"
	"github.com/HongJungWan/commerce-system/test/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestOrderInteractor_CreateOrder_Failure_InvalidMember(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	orderRepo := repository.NewOrderRepository(db)
	memberRepo := repository.NewMemberRepository(db)
	productRepo := repository.NewProductRepository(db)
	interactor := usecases.NewOrderInteractor(orderRepo, memberRepo, productRepo)

	product := &domain.Product{
		ProductNumber: "P12345",
		ProductName:   "Test Product",
		Price:         1000,
		StockQuantity: 10,
	}
	_ = productRepo.Create(product)

	req := &request.CreateOrderRequest{
		OrderNumber:   "O12345",
		ProductNumber: "P12345",
		Quantity:      2,
	}

	// When
	responseData, err := interactor.CreateOrder(req, "InvalidMember")

	// Then
	assert.Error(t, err)
	assert.Nil(t, responseData)
	assert.Equal(t, "유효하지 않은 회원 번호입니다.", err.Error())
}

func TestOrderInteractor_CreateOrder_Failure_InvalidProduct(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	orderRepo := repository.NewOrderRepository(db)
	memberRepo := repository.NewMemberRepository(db)
	productRepo := repository.NewProductRepository(db)
	interactor := usecases.NewOrderInteractor(orderRepo, memberRepo, productRepo)

	member := &domain.Member{
		MemberNumber: "M12345",
		Username:     "testuser",
		FullName:     "Test User",
		Email:        "testuser@example.com",
	}
	member.SetPassword("password123")
	_ = memberRepo.Create(member)

	req := &request.CreateOrderRequest{
		OrderNumber:   "O12345",
		ProductNumber: "InvalidProduct",
		Quantity:      2,
	}

	// When
	responseData, err := interactor.CreateOrder(req, "M12345")

	// Then
	assert.Error(t, err)
	assert.Nil(t, responseData)
	assert.Equal(t, "유효하지 않은 상품 번호입니다.", err.Error())
}

func TestOrderInteractor_CreateOrder_Failure_InsufficientStock(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	orderRepo := repository.NewOrderRepository(db)
	memberRepo := repository.NewMemberRepository(db)
	productRepo := repository.NewProductRepository(db)
	interactor := usecases.NewOrderInteractor(orderRepo, memberRepo, productRepo)

	member := &domain.Member{
		MemberNumber: "M12345",
		Username:     "testuser",
		FullName:     "Test User",
		Email:        "testuser@example.com",
	}
	member.SetPassword("password123")
	_ = memberRepo.Create(member)

	product := &domain.Product{
		ProductNumber: "P12345",
		ProductName:   "Test Product",
		Price:         1000,
		StockQuantity: 1,
	}
	_ = productRepo.Create(product)

	req := &request.CreateOrderRequest{
		OrderNumber:   "O12345",
		ProductNumber: "P12345",
		Quantity:      2,
	}

	// When
	responseData, err := interactor.CreateOrder(req, "M12345")

	// Then
	assert.Error(t, err)
	assert.Nil(t, responseData)
	assert.Equal(t, "재고 수량이 부족합니다.", err.Error())
}

func TestOrderInteractor_CancelOrder_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	orderRepo := repository.NewOrderRepository(db)
	memberRepo := repository.NewMemberRepository(db)
	productRepo := repository.NewProductRepository(db)
	interactor := usecases.NewOrderInteractor(orderRepo, memberRepo, productRepo)

	member := &domain.Member{
		MemberNumber: "M12345",
		Username:     "testuser",
		FullName:     "Test User",
		Email:        "testuser@example.com",
	}
	member.SetPassword("password123")
	_ = memberRepo.Create(member)

	product := &domain.Product{
		ProductNumber: "P12345",
		ProductName:   "Test Product",
		Price:         1000,
		StockQuantity: 10,
	}
	_ = productRepo.Create(product)

	order := &domain.Order{
		OrderNumber:   "O12345",
		OrderDate:     time.Now(),
		MemberNumber:  "M12345",
		ProductNumber: "P12345",
		Price:         1000,
		Quantity:      2,
		TotalAmount:   2000,
		IsCanceled:    false,
	}
	_ = orderRepo.Create(order)

	// When
	err := interactor.CancelOrder("O12345", "M12345")

	// Then
	assert.NoError(t, err)

	updatedOrder, _ := orderRepo.GetByOrderNumber("O12345")
	assert.True(t, updatedOrder.IsCanceled)

	updatedProduct, _ := productRepo.GetByProductNumber("P12345")
	assert.Equal(t, 12, updatedProduct.StockQuantity)
}

func TestOrderInteractor_CancelOrder_Failure_OrderNotFound(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	orderRepo := repository.NewOrderRepository(db)
	memberRepo := repository.NewMemberRepository(db)
	productRepo := repository.NewProductRepository(db)
	interactor := usecases.NewOrderInteractor(orderRepo, memberRepo, productRepo)

	// When
	err := interactor.CancelOrder("InvalidOrder", "M12345")

	// Then
	assert.Error(t, err)
}

func TestOrderInteractor_CancelOrder_Failure_Unauthorized(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	orderRepo := repository.NewOrderRepository(db)
	memberRepo := repository.NewMemberRepository(db)
	productRepo := repository.NewProductRepository(db)
	interactor := usecases.NewOrderInteractor(orderRepo, memberRepo, productRepo)

	member := &domain.Member{
		MemberNumber: "M12345",
		Username:     "testuser",
		FullName:     "Test User",
		Email:        "testuser@example.com",
	}
	member.SetPassword("password123")
	_ = memberRepo.Create(member)

	order := &domain.Order{
		OrderNumber:   "O12345",
		OrderDate:     time.Now(),
		MemberNumber:  "M99999",
		ProductNumber: "P12345",
		Price:         1000,
		Quantity:      2,
		TotalAmount:   2000,
		IsCanceled:    false,
	}
	_ = orderRepo.Create(order)

	// When
	err := interactor.CancelOrder("O12345", "M12345")

	// Then
	assert.Error(t, err)
	assert.Equal(t, "해당 주문에 대한 권한이 없습니다.", err.Error())
}

func TestOrderInteractor_GetMyOrders_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	orderRepo := repository.NewOrderRepository(db)
	memberRepo := repository.NewMemberRepository(db)
	productRepo := repository.NewProductRepository(db)
	interactor := usecases.NewOrderInteractor(orderRepo, memberRepo, productRepo)

	member := &domain.Member{
		MemberNumber: "M12345",
		Username:     "testuser",
		FullName:     "Test User",
		Email:        "testuser@example.com",
	}
	member.SetPassword("password123")
	_ = memberRepo.Create(member)

	order1 := &domain.Order{
		OrderNumber:   "O12345",
		OrderDate:     time.Now(),
		MemberNumber:  "M12345",
		ProductNumber: "P12345",
		Price:         1000,
		Quantity:      2,
		TotalAmount:   2000,
		IsCanceled:    false,
	}
	order2 := &domain.Order{
		OrderNumber:   "O12346",
		OrderDate:     time.Now(),
		MemberNumber:  "M12345",
		ProductNumber: "P12346",
		Price:         1500,
		Quantity:      1,
		TotalAmount:   1500,
		IsCanceled:    false,
	}
	_ = orderRepo.Create(order1)
	_ = orderRepo.Create(order2)

	// When
	orders, err := interactor.GetMyOrders("M12345")

	// Then
	assert.NoError(t, err)
	assert.Len(t, orders, 2)
}

func TestOrderInteractor_GetMonthlyStats_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	orderRepo := repository.NewOrderRepository(db)
	memberRepo := repository.NewMemberRepository(db)
	productRepo := repository.NewProductRepository(db)
	interactor := usecases.NewOrderInteractor(orderRepo, memberRepo, productRepo)

	order1 := &domain.Order{
		OrderNumber:   "O12345",
		OrderDate:     time.Date(2024, 9, 10, 0, 0, 0, 0, time.UTC),
		MemberNumber:  "M12345",
		ProductNumber: "P12345",
		Price:         1000,
		Quantity:      2,
		TotalAmount:   2000,
		IsCanceled:    false,
	}
	order2 := &domain.Order{
		OrderNumber:   "O12346",
		OrderDate:     time.Date(2024, 9, 15, 0, 0, 0, 0, time.UTC),
		MemberNumber:  "M12346",
		ProductNumber: "P12346",
		Price:         1500,
		Quantity:      1,
		TotalAmount:   1500,
		IsCanceled:    true,
	}
	_ = orderRepo.Create(order1)
	_ = orderRepo.Create(order2)

	// When
	totalSales, totalCanceled, err := interactor.GetMonthlyStats("2024-09")

	// Then
	assert.NoError(t, err)
	assert.EqualValues(t, 2000, totalSales)
	assert.EqualValues(t, 1500, totalCanceled)
}

func TestOrderInteractor_GetMonthlyStats_Failure_InvalidMonth(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	orderRepo := repository.NewOrderRepository(db)
	memberRepo := repository.NewMemberRepository(db)
	productRepo := repository.NewProductRepository(db)
	interactor := usecases.NewOrderInteractor(orderRepo, memberRepo, productRepo)

	// When
	totalSales, totalCanceled, err := interactor.GetMonthlyStats("invalid-month")

	// Then
	assert.Error(t, err)
	assert.EqualValues(t, 0, totalSales)
	assert.EqualValues(t, 0, totalCanceled)
}
