package usecases_test

import (
	"testing"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/infrastructure/repository"
	"github.com/HongJungWan/commerce-system/internal/interfaces/dto/request"
	"github.com/HongJungWan/commerce-system/internal/usecases"
	"github.com/HongJungWan/commerce-system/test/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestProductInteractor_CreateProduct_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	productRepo := repository.NewProductRepository(db)
	interactor := usecases.NewProductInteractor(productRepo, db)

	req := &request.CreateProductRequest{
		ProductName:   "New Product",
		Category:      "Electronics",
		Price:         1000,
		StockQuantity: 10,
	}

	// When
	responseData, err := interactor.CreateProduct(req)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "상품이 등록되었습니다.", responseData.Message)
	assert.Equal(t, "New Product", responseData.Product.ProductName)

	retrievedProduct, err := productRepo.GetById(responseData.Product.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedProduct)
	assert.Equal(t, "New Product", retrievedProduct.ProductName)
}

func TestProductInteractor_CreateProduct_Failure_InvalidProduct(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	productRepo := repository.NewProductRepository(db)
	interactor := usecases.NewProductInteractor(productRepo, db)

	req := &request.CreateProductRequest{
		ProductName:   "", // 상품명 누락
		Category:      "Electronics",
		Price:         -1000, // 잘못된 가격
		StockQuantity: -10,   // 잘못된 재고 수량
	}

	// When
	_, err := interactor.CreateProduct(req)

	// Then
	assert.Error(t, err)
	assert.Equal(t, "상품명이 누락되었습니다.", err.Error())
}

func TestProductInteractor_GetProducts_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	productRepo := repository.NewProductRepository(db)
	interactor := usecases.NewProductInteractor(productRepo, db)

	product1 := &domain.Product{
		ProductNumber: "P12345",
		ProductName:   "Product One",
		Category:      "Electronics",
		Price:         1000,
		StockQuantity: 10,
	}
	product2 := &domain.Product{
		ProductNumber: "P12346",
		ProductName:   "Product Two",
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
	assert.Equal(t, "Product One", products[0].ProductName)
}

func TestProductInteractor_UpdateStock_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	productRepo := repository.NewProductRepository(db)
	interactor := usecases.NewProductInteractor(productRepo, db)

	product := &domain.Product{
		ProductNumber: "P12345",
		ProductName:   "Test Product",
		Price:         1000,
		StockQuantity: 10,
	}
	_ = productRepo.Create(product)

	// When
	err := interactor.UpdateStock(product.ID, 20)

	// Then
	assert.NoError(t, err)
	updatedProduct, _ := productRepo.GetById(product.ID)
	assert.Equal(t, 20, updatedProduct.StockQuantity)
}

func TestProductInteractor_UpdateStock_Failure_ProductNotFound(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	productRepo := repository.NewProductRepository(db)
	interactor := usecases.NewProductInteractor(productRepo, db)

	// When
	err := interactor.UpdateStock(9999, 20) // 존재하지 않는 ID

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
		ProductName:   "Test Product",
		Price:         1000,
		StockQuantity: 10,
	}
	_ = productRepo.Create(product)

	// When
	err := interactor.DeleteProduct(product.ID)

	// Then
	assert.NoError(t, err)
	deletedProduct, err := productRepo.GetById(product.ID)
	assert.Error(t, err)
	assert.Nil(t, deletedProduct)
}

func TestProductInteractor_DeleteProduct_Failure_HasOrders(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db) // OrderRepository가 정의되어 있다고 가정
	interactor := usecases.NewProductInteractor(productRepo, db)

	product := &domain.Product{
		ProductNumber: "P12345",
		ProductName:   "Test Product",
		Price:         1000,
		StockQuantity: 10,
	}
	_ = productRepo.Create(product)

	order := &domain.Order{
		OrderNumber:   "O12345",
		ProductNumber: product.ProductNumber,
	}
	_ = orderRepo.Create(order)

	// When
	err := interactor.DeleteProduct(product.ID)

	// Then
	assert.Error(t, err)
	assert.Equal(t, "주문된 이력이 있어 삭제할 수 없습니다.", err.Error())
}
