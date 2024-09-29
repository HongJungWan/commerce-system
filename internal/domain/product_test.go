package domain_test

import (
	"github.com/HongJungWan/commerce-system/test/fixtures"
	"testing"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestProduct_Validate_Success(t *testing.T) {
	// Given
	product := &domain.Product{
		ProductNumber: "P12345",
		Name:          "Test Product",
		Price:         1000,
		StockQuantity: 10,
	}

	// When
	err := product.Validate()

	// Then
	assert.NoError(t, err)
}

func TestProduct_Validate_Failure_MissingFields(t *testing.T) {
	// Given
	product := &domain.Product{
		ProductNumber: "",
		Name:          "",
		Price:         -1000,
		StockQuantity: -10,
	}

	// When
	err := product.Validate()

	// Then
	assert.Error(t, err)
	assert.Equal(t, "필수 필드가 누락되었거나 잘못된 값입니다.", err.Error())
}

func TestProduct_UpdateStock_Success(t *testing.T) {
	// Given
	product := &domain.Product{
		StockQuantity: 10,
	}

	// When
	err := product.UpdateStock(20)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 20, product.StockQuantity)
}

func TestProduct_UpdateStock_Failure_NegativeQuantity(t *testing.T) {
	// Given
	product := &domain.Product{
		StockQuantity: 10,
	}

	// When
	err := product.UpdateStock(-5)

	// Then
	assert.Error(t, err)
	assert.Equal(t, "재고 수량은 음수일 수 없습니다.", err.Error())
}

func TestProduct_CanBeDeleted_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	product := &domain.Product{
		ProductNumber: "P12345",
	}
	db.Create(product)

	// When
	canBeDeleted, err := product.CanBeDeleted(db)

	// Then
	assert.NoError(t, err)
	assert.True(t, canBeDeleted)
}

func TestProduct_CanBeDeleted_Failure_HasOrders(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	product := &domain.Product{
		ProductNumber: "P12345",
	}
	order := &domain.Order{
		OrderNumber:   "O12345",
		ProductNumber: "P12345",
	}

	db.Create(product)
	db.Create(order)

	// When
	canBeDeleted, err := product.CanBeDeleted(db)

	// Then
	assert.NoError(t, err)
	assert.False(t, canBeDeleted)
}
