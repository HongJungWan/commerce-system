package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/infrastructure/repository"
	"github.com/HongJungWan/commerce-system/internal/interfaces/controller"
	"github.com/HongJungWan/commerce-system/internal/usecases"
	"github.com/HongJungWan/commerce-system/test/fixtures"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestProductController_GetProducts_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	productRepo := repository.NewProductRepository(db)
	productInteractor := usecases.NewProductInteractor(productRepo, db)
	productController := controller.NewProductController(productInteractor)

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

	router := gin.Default()
	router.GET("/products", productController.GetProducts)

	req, _ := http.NewRequest("GET", "/products?category=Electronics", nil)

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusOK, resp.Code)
	var products []domain.Product
	err := json.Unmarshal(resp.Body.Bytes(), &products)
	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, "Product One", products[0].ProductName)
}

func TestProductController_CreateProduct_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	productRepo := repository.NewProductRepository(db)
	productInteractor := usecases.NewProductInteractor(productRepo, db)
	productController := controller.NewProductController(productInteractor)

	router := gin.Default()
	router.POST("/products", func(c *gin.Context) {
		c.Set("is_admin", true)
		productController.CreateProduct(c)
	})

	newProduct := domain.Product{
		ProductNumber: "P12345",
		ProductName:   "New Product",
		Price:         1000,
		StockQuantity: 10,
	}
	requestBody, _ := json.Marshal(newProduct)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusCreated, resp.Code)
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "상품이 등록되었습니다.", response["message"])
}

func TestProductController_CreateProduct_Failure_Unauthorized(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	productRepo := repository.NewProductRepository(db)
	productInteractor := usecases.NewProductInteractor(productRepo, db)
	productController := controller.NewProductController(productInteractor)

	router := gin.Default()
	router.POST("/products", func(c *gin.Context) {
		c.Set("is_admin", false)
		productController.CreateProduct(c)
	})

	newProduct := domain.Product{
		ProductNumber: "P12345",
		ProductName:   "New Product",
		Price:         1000,
		StockQuantity: 10,
	}
	requestBody, _ := json.Marshal(newProduct)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusForbidden, resp.Code)
}

func TestProductController_UpdateStock_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	productRepo := repository.NewProductRepository(db)
	productInteractor := usecases.NewProductInteractor(productRepo, db)
	productController := controller.NewProductController(productInteractor)

	product := &domain.Product{
		ProductNumber: "P12345",
		ProductName:   "Test Product",
		Price:         1000,
		StockQuantity: 10,
	}
	_ = productRepo.Create(product)

	router := gin.Default()
	router.PUT("/products/:product_number/stock", func(c *gin.Context) {
		c.Set("is_admin", true)
		productController.UpdateStock(c)
	})

	updateData := map[string]interface{}{
		"stock_quantity": 20,
	}
	requestBody, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", "/products/P12345/stock", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusOK, resp.Code)
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "재고 수량이 수정되었습니다.", response["message"])

	updatedProduct, _ := productRepo.GetById("P12345")
	assert.Equal(t, 20, updatedProduct.StockQuantity)
}

func TestProductController_DeleteProduct_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	productRepo := repository.NewProductRepository(db)
	productInteractor := usecases.NewProductInteractor(productRepo, db)
	productController := controller.NewProductController(productInteractor)

	product := &domain.Product{
		ProductNumber: "P12345",
		ProductName:   "Test Product",
		Price:         1000,
		StockQuantity: 10,
	}
	_ = productRepo.Create(product)

	router := gin.Default()
	router.DELETE("/products/:product_number", func(c *gin.Context) {
		c.Set("is_admin", true)
		productController.DeleteProduct(c)
	})

	req, _ := http.NewRequest("DELETE", "/products/P12345", nil)

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusOK, resp.Code)
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "상품이 삭제되었습니다.", response["message"])

	deletedProduct, err := productRepo.GetById("P12345")
	assert.Error(t, err)
	assert.Nil(t, deletedProduct)
}

func TestProductController_DeleteProduct_Failure_HasOrders(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	productInteractor := usecases.NewProductInteractor(productRepo, db)
	productController := controller.NewProductController(productInteractor)

	product := &domain.Product{
		ProductNumber: "P12345",
		ProductName:   "Test Product",
		Price:         1000,
		StockQuantity: 10,
	}
	order := &domain.Order{
		OrderNumber:   "O12345",
		ProductNumber: "P12345",
	}
	_ = productRepo.Create(product)
	_ = orderRepo.Create(order)

	router := gin.Default()
	router.DELETE("/products/:product_number", func(c *gin.Context) {
		c.Set("is_admin", true)
		productController.DeleteProduct(c)
	})

	req, _ := http.NewRequest("DELETE", "/products/P12345", nil)

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "주문된 이력이 있어 삭제할 수 없습니다.", response["error"])
}
