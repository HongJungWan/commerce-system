package usecases_test

import (
	"testing"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/infrastructure/repository"
	"github.com/HongJungWan/commerce-system/internal/usecases"
	"github.com/HongJungWan/commerce-system/test/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestProductInteractor_CreateProduct_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	productRepo := repository.NewProductRepository(db)
	interactor := usecases.NewProductInteractor(productRepo, db)

	product := &domain.Product{
		ProductNumber: "P12345",
		Name:          "New Product",
		Price:         1000,
		StockQuantity: 10,
	}

	// When
	err := interactor.CreateProduct(product)

	// Then
	assert.NoError(t, err)
	retrievedProduct, _ := productRepo.GetByProductNumber("P12345")
	assert.NotNil(t, retrievedProduct)
}

func TestProductInteractor_CreateProduct_Failure_InvalidProduct(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	productRepo := repository.NewProductRepository(db)
	interactor := usecases.NewProductInteractor(productRepo, db)

	product := &domain.Product{
		ProductNumber: "",
		Name:          "",
		Price:         -1000,
		StockQuantity: -10,
	}

	// When
	err := interactor.CreateProduct(product)

	// Then
	assert.Error(t, err)
	assert.Equal(t, "필수 필드가 누락되었거나 잘못된 값입니다.", err.Error())
}

func TestProductInteractor_GetProducts_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	productRepo := repository.NewProductRepository(db)
	interactor := usecases.NewProductInteractor(productRepo, db)

	product1 := &domain.Product{
		ProductNumber: "P12345",
		Name:          "Product One",
		Category:      "Electronics",
		Price:         1000,
		StockQuantity: 10,
	}
	product2 := &domain.Product{
		ProductNumber: "P12346",
		Name:          "Product Two",
		Category:      "Home",
		Price:         1500,
		StockQuantity: 5,
	}
	_ = productRepo.Create(product1)
	_ = productRepo.Create(product2)

	// When
	products, err := interactor.GetProducts(map[string]interface{}{
		"category": "Electronics",
	})

	// Then
	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, "Product One", products[0].Name)
}

func TestProductInteractor_UpdateStock_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	productRepo := repository.NewProductRepository(db)
	interactor := usecases.NewProductInteractor(productRepo, db)

	product := &domain.Product{
		ProductNumber: "P12345",
		Name:          "Test Product",
		Price:         1000,
		StockQuantity: 10,
	}
	_ = productRepo.Create(product)

	// When
	err := interactor.UpdateStock("P12345", 20)

	// Then
	assert.NoError(t, err)
	updatedProduct, _ := productRepo.GetByProductNumber("P12345")
	assert.Equal(t, 20, updatedProduct.StockQuantity)
}

func TestProductInteractor_UpdateStock_Failure_ProductNotFound(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	productRepo := repository.NewProductRepository(db)
	interactor := usecases.NewProductInteractor(productRepo, db)

	// When
	err := interactor.UpdateStock("InvalidProduct", 20)

	// Then
	assert.Error(t, err)
}

func TestProductInteractor_DeleteProduct_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	productRepo := repository.NewProductRepository(db)
	interactor := usecases.NewProductInteractor(productRepo, db)

	product := &domain.Product{
		ProductNumber: "P12345",
		Name:          "Test Product",
		Price:         1000,
		StockQuantity: 10,
	}
	_ = productRepo.Create(product)

	// When
	err := interactor.DeleteProduct("P12345")

	// Then
	assert.NoError(t, err)
	deletedProduct, err := productRepo.GetByProductNumber("P12345")
	assert.Error(t, err)
	assert.Nil(t, deletedProduct)
}

func TestProductInteractor_DeleteProduct_Failure_HasOrders(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	interactor := usecases.NewProductInteractor(productRepo, db)

	product := &domain.Product{
		ProductNumber: "P12345",
		Name:          "Test Product",
		Price:         1000,
		StockQuantity: 10,
	}
	order := &domain.Order{
		OrderNumber:   "O12345",
		ProductNumber: "P12345",
	}
	_ = productRepo.Create(product)
	_ = orderRepo.Create(order)

	// When
	err := interactor.DeleteProduct("P12345")

	// Then
	assert.Error(t, err)
	assert.Equal(t, "주문된 이력이 있어 삭제할 수 없습니다.", err.Error())
}
