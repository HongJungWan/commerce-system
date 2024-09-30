package controller

import (
	"github.com/HongJungWan/commerce-system/internal/domain"
	"net/http"

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
	c.JSON(http.StatusOK, products)
}

func (pc *ProductController) CreateProduct(c *gin.Context) {
	isAdmin := c.GetBool("is_admin")
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "접근 권한이 없습니다."})
		return
	}

	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다."})
		return
	}

	if err := pc.productInteractor.CreateProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "상품이 등록되었습니다.", "product": product})
}

func (pc *ProductController) UpdateStock(c *gin.Context) {
	isAdmin := c.GetBool("is_admin")
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "접근 권한이 없습니다."})
		return
	}

	productNumber := c.Param("product_number")
	var req struct {
		StockQuantity int `json:"stock_quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다."})
		return
	}

	if err := pc.productInteractor.UpdateStock(productNumber, req.StockQuantity); err != nil {
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

	productNumber := c.Param("product_number")
	if err := pc.productInteractor.DeleteProduct(productNumber); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "상품이 삭제되었습니다."})
}
