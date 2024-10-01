package controller

import (
	"net/http"
	"strconv"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/interfaces/dto/request"
	"github.com/HongJungWan/commerce-system/internal/interfaces/dto/response"
	"github.com/HongJungWan/commerce-system/internal/usecases"
	"github.com/gin-gonic/gin"
)

type ProductController struct {
	productInteractor *usecases.ProductInteractor
}

func NewProductController(pi *usecases.ProductInteractor) *ProductController {
	return &ProductController{productInteractor: pi}
}

func (pc *ProductController) GetProducts(c *gin.Context) {
	filter := make(map[string]interface{})
	if category := c.Query("category"); category != "" {
		filter["category"] = category
	}
	if name := c.Query("product_name"); name != "" {
		filter["product_name"] = name
	}

	products, err := pc.productInteractor.GetProducts(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "상품 목록을 가져올 수 없습니다."})
		return
	}

	var productResponses []response.ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, response.ProductResponse{
			ID:            product.ID,
			ProductNumber: product.ProductNumber,
			ProductName:   product.ProductName,
			Category:      product.Category,
			Price:         product.Price,
			StockQuantity: product.StockQuantity,
		})
	}

	c.JSON(http.StatusOK, productResponses)
}

func (pc *ProductController) CreateProduct(c *gin.Context) {
	isAdmin := c.GetBool("is_admin")
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "접근 권한이 없습니다."})
		return
	}

	var req request.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다."})
		return
	}

	product := &domain.Product{
		ProductNumber: req.ProductNumber,
		ProductName:   req.ProductName,
		Category:      req.Category,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
	}

	if err := pc.productInteractor.CreateProduct(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	productResponse := response.ProductResponse{
		ID:            product.ID,
		ProductNumber: product.ProductNumber,
		ProductName:   product.ProductName,
		Category:      product.Category,
		Price:         product.Price,
		StockQuantity: product.StockQuantity,
	}

	c.JSON(http.StatusCreated, response.CreateProductResponse{
		Message: "상품이 등록되었습니다.",
		Product: productResponse,
	})
}

func (pc *ProductController) UpdateStock(c *gin.Context) {
	isAdmin := c.GetBool("is_admin")
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "접근 권한이 없습니다."})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 상품 ID입니다."})
		return
	}

	var req request.UpdateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다."})
		return
	}

	if err := pc.productInteractor.UpdateStock(id, req.StockQuantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "재고 수량이 수정되었습니다."})
}

func (pc *ProductController) DeleteProduct(c *gin.Context) {
	isAdmin := c.GetBool("is_admin")
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "접근 권한이 없습니다."})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 상품 ID입니다."})
		return
	}

	if err := pc.productInteractor.DeleteProduct(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "상품이 삭제되었습니다."})
}
