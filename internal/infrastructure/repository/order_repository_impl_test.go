package repository_test

import (
	"testing"
	"time"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/infrastructure/repository"
	"github.com/HongJungWan/commerce-system/test/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestOrderRepositoryImpl_Create_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	repo := repository.NewOrderRepository(db)
	order := &domain.Order{
		OrderNumber:   "O12345",
		OrderDate:     time.Now(),
		MemberNumber:  "M12345",
		ProductNumber: 12345,
		Price:         1000,
		Quantity:      2,
		TotalAmount:   2000,
		IsCanceled:    true,
	}

	// When
	err := repo.Create(order)

	// Then
	assert.NoError(t, err)
}

func TestOrderRepositoryImpl_Create_Failure_DuplicateOrderNumber(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	repo := repository.NewOrderRepository(db)
	order1 := &domain.Order{
		OrderNumber:   "O12345",
		OrderDate:     time.Now(),
		MemberNumber:  "M12345",
		ProductNumber: 12345,
		Price:         1000,
		Quantity:      2,
		TotalAmount:   2000,
		IsCanceled:    true,
	}
	order2 := &domain.Order{
		OrderNumber:   "O12345",
		OrderDate:     time.Now(),
		MemberNumber:  "M12346",
		ProductNumber: 12346,
		Price:         2000,
		Quantity:      1,
		TotalAmount:   2000,
		IsCanceled:    true,
	}
	_ = repo.Create(order1)

	// When
	err := repo.Create(order2)

	// Then
	assert.Error(t, err)
}

func TestOrderRepositoryImpl_GetByOrderNumber_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	repo := repository.NewOrderRepository(db)
	order := &domain.Order{
		OrderNumber:   "O12345",
		OrderDate:     time.Now(),
		MemberNumber:  "M12345",
		ProductNumber: 12345,
		Price:         1000,
		Quantity:      2,
		TotalAmount:   2000,
		IsCanceled:    true,
	}
	_ = repo.Create(order)

	// When
	retrievedOrder, err := repo.GetByOrderNumber("O12345")

	// Then
	assert.NoError(t, err)
	assert.Equal(t, order.MemberNumber, retrievedOrder.MemberNumber)
}

func TestOrderRepositoryImpl_GetByOrderNumber_Failure_NotFound(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	repo := repository.NewOrderRepository(db)

	// When
	retrievedOrder, err := repo.GetByOrderNumber("nonexistent")

	// Then
	assert.Error(t, err)
	assert.Nil(t, retrievedOrder)
}

func TestOrderRepositoryImpl_GetByMemberNumber_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	repo := repository.NewOrderRepository(db)
	order1 := &domain.Order{
		OrderNumber:   "O12345",
		OrderDate:     time.Now(),
		MemberNumber:  "M12345",
		ProductNumber: 12345,
		Price:         1000,
		Quantity:      2,
		TotalAmount:   2000,
		IsCanceled:    true,
	}
	order2 := &domain.Order{
		OrderNumber:   "O12346",
		OrderDate:     time.Now(),
		MemberNumber:  "M12345",
		ProductNumber: 12346,
		Price:         1500,
		Quantity:      1,
		TotalAmount:   1500,
		IsCanceled:    true,
	}
	_ = repo.Create(order1)
	_ = repo.Create(order2)

	// When
	orders, err := repo.GetByMemberNumber("M12345")

	// Then
	assert.NoError(t, err)
	assert.Len(t, orders, 2)
}

func TestOrderRepositoryImpl_Update_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	repo := repository.NewOrderRepository(db)
	order := &domain.Order{
		OrderNumber:   "O12345",
		OrderDate:     time.Now(),
		MemberNumber:  "M12345",
		ProductNumber: 12345,
		Price:         1000,
		Quantity:      2,
		TotalAmount:   2000,
		IsCanceled:    true,
	}
	_ = repo.Create(order)

	// When
	order.IsCanceled = true
	err := repo.Update(order)

	// Then
	assert.NoError(t, err)

	// Verify
	updatedOrder, _ := repo.GetByOrderNumber("O12345")
	assert.Equal(t, true, updatedOrder.IsCanceled)
}

func TestOrderRepositoryImpl_GetMonthlyStats_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	repo := repository.NewOrderRepository(db)
	order1 := &domain.Order{
		OrderNumber:   "O12345",
		OrderDate:     time.Date(2024, 9, 10, 0, 0, 0, 0, time.UTC),
		MemberNumber:  "M12345",
		ProductNumber: 12345,
		Price:         1000,
		Quantity:      2,
		TotalAmount:   2000,
		IsCanceled:    true,
	}
	order2 := &domain.Order{
		OrderNumber:   "O12346",
		OrderDate:     time.Date(2024, 9, 15, 0, 0, 0, 0, time.UTC),
		MemberNumber:  "M12346",
		ProductNumber: 12346,
		Price:         1500,
		Quantity:      1,
		TotalAmount:   1500,
		IsCanceled:    true,
	}
	_ = repo.Create(order1)
	_ = repo.Create(order2)

	// When
	totalSales, totalCanceled, err := repo.GetMonthlyStats("2024-09")

	// Then
	assert.NoError(t, err)
	assert.EqualValues(t, 2000, totalSales)
	assert.EqualValues(t, 1500, totalCanceled)
}

func TestOrderRepositoryImpl_GetMonthlyStats_Failure_InvalidMonth(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	repo := repository.NewOrderRepository(db)

	// When
	totalSales, totalCanceled, err := repo.GetMonthlyStats("invalid-month")

	// Then
	assert.Error(t, err)
	assert.EqualValues(t, 0, totalSales)
	assert.EqualValues(t, 0, totalCanceled)
}
