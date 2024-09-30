package domain_test

import (
	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrder_Validate_Success(t *testing.T) {
	// Given
	order := &domain.Order{
		OrderNumber:   "O12345",
		MemberNumber:  "M12345",
		ProductNumber: 12345,
		Quantity:      2,
		Price:         1000,
	}

	// When
	err := order.Validate()

	// Then
	assert.NoError(t, err)
}

func TestOrder_Validate_Failure_MissingFields(t *testing.T) {
	// Given
	order := &domain.Order{
		OrderNumber:   "",
		MemberNumber:  "",
		ProductNumber: 0,
		Quantity:      -1,
		Price:         -1000,
	}

	// When
	err := order.Validate()

	// Then
	assert.Error(t, err)
	assert.Equal(t, "주문번호가 누락되었습니다.", err.Error())
}

func TestOrder_Cancel_Success(t *testing.T) {
	// Given
	order := &domain.Order{
		OrderNumber: "O12345",
		IsCanceled:  true,
	}

	// When
	err := order.Cancel()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, true, order.IsCanceled)
	assert.NotNil(t, order.CanceledAt)
}

func TestOrder_Cancel_Failure_AlreadyCanceled(t *testing.T) {
	// Given
	order := &domain.Order{
		OrderNumber: "O12345",
		IsCanceled:  true,
	}

	// When
	err := order.Cancel()

	// Then
	assert.Error(t, err)
	assert.Equal(t, "이미 취소된 주문입니다.", err.Error())
}
