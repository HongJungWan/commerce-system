package repository_test

import (
	"testing"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/infrastructure/repository"
	"github.com/HongJungWan/commerce-system/test/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestProductRepositoryImpl_Create_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	repo := repository.NewProductRepository(db)
	product := &domain.Product{
		ProductNumber: "P12345",
		Name:          "Test Product",
		Category:      "Electronics",
		Price:         1000,
		StockQuantity: 10,
	}

	// When
	err := repo.Create(product)

	// Then
	assert.NoError(t, err)
}

func TestProductRepositoryImpl_Create_Failure_DuplicateProductNumber(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	repo := repository.NewProductRepository(db)
	product1 := &domain.Product{
		ProductNumber: "P12345",
		Name:          "Product One",
		Category:      "Electronics",
		Price:         1000,
		StockQuantity: 10,
	}
	product2 := &domain.Product{
		ProductNumber: "P12345",
		Name:          "Product Two",
		Category:      "Home",
		Price:         1500,
		StockQuantity: 5,
	}
	_ = repo.Create(product1)

	// When
	err := repo.Create(product2)

	// Then
	assert.Error(t, err)
}

func TestProductRepositoryImpl_GetAll_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	repo := repository.NewProductRepository(db)
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
	_ = repo.Create(product1)
	_ = repo.Create(product2)

	// When
	products, err := repo.GetAll(map[string]interface{}{})

	// Then
	assert.NoError(t, err)
	assert.Len(t, products, 2)
}

func TestProductRepositoryImpl_GetAll_WithFilters(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	repo := repository.NewProductRepository(db)
	product1 := &domain.Product{
		ProductNumber: "P12345",
		Name:          "Smartphone",
		Category:      "Electronics",
		Price:         1000,
		StockQuantity: 10,
	}
	product2 := &domain.Product{
		ProductNumber: "P12346",
		Name:          "Vacuum Cleaner",
		Category:      "Home",
		Price:         1500,
		StockQuantity: 5,
	}
	_ = repo.Create(product1)
	_ = repo.Create(product2)

	// When
	products, err := repo.GetAll(map[string]interface{}{
		"category": "Electronics",
		"name":     "Smart",
	})

	// Then
	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, "Smartphone", products[0].Name)
}

func TestProductRepositoryImpl_GetByProductNumber_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	repo := repository.NewProductRepository(db)
	product := &domain.Product{
		ProductNumber: "P12345",
		Name:          "Test Product",
		Category:      "Electronics",
		Price:         1000,
		StockQuantity: 10,
	}
	_ = repo.Create(product)

	// When
	retrievedProduct, err := repo.GetByProductNumber("P12345")

	// Then
	assert.NoError(t, err)
	assert.Equal(t, product.Name, retrievedProduct.Name)
}

func TestProductRepositoryImpl_GetByProductNumber_Failure_NotFound(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	repo := repository.NewProductRepository(db)

	// When
	retrievedProduct, err := repo.GetByProductNumber("nonexistent")

	// Then
	assert.Error(t, err)
	assert.Nil(t, retrievedProduct)
}

func TestProductRepositoryImpl_Update_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	repo := repository.NewProductRepository(db)
	product := &domain.Product{
		ProductNumber: "P12345",
		Name:          "Old Name",
		Category:      "Electronics",
		Price:         1000,
		StockQuantity: 10,
	}
	_ = repo.Create(product)

	// When
	product.Name = "New Name"
	err := repo.Update(product)

	// Then
	assert.NoError(t, err)

	// Verify
	updatedProduct, _ := repo.GetByProductNumber("P12345")
	assert.Equal(t, "New Name", updatedProduct.Name)
}

func TestProductRepositoryImpl_Delete_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	repo := repository.NewProductRepository(db)
	product := &domain.Product{
		ProductNumber: "P12345",
		Name:          "Test Product",
		Category:      "Electronics",
		Price:         1000,
		StockQuantity: 10,
	}
	_ = repo.Create(product)

	// When
	err := repo.Delete("P12345")

	// Then
	assert.NoError(t, err)

	// Verify
	deletedProduct, err := repo.GetByProductNumber("P12345")
	assert.Error(t, err)
	assert.Nil(t, deletedProduct)
}
