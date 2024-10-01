package controller

import (
	"net/http"
	"strconv"

	"github.com/HongJungWan/commerce-system/internal/interfaces/dto/request"
	"github.com/HongJungWan/commerce-system/internal/usecases"
	"github.com/gin-gonic/gin"
)

type ProductController struct {
	productInteractor *usecases.ProductInteractor
}

func NewProductController(pi *usecases.ProductInteractor) *ProductController {
	return &ProductController{productInteractor: pi}
}

// GetProducts godoc
// @Summary      상품 목록 조회
// @Description  상품 목록을 필터링하여 조회합니다.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        category query string false "카테고리"
// @Param        product_name query string false "상품명"
// @Success      200 {array} response.ProductResponse "상품 목록"
// @Failure      500 {object} map[string]string "조회 실패"
// @Router       /products [get]
func (pc *ProductController) GetProducts(c *gin.Context) {
	filter := make(map[string]interface{})
	if category := c.Query("category"); category != "" {
		filter["category"] = category
	}
	if name := c.Query("product_name"); name != "" {
		filter["product_name"] = name
	}

	productResponses, err := pc.productInteractor.GetProducts(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "상품 목록을 가져올 수 없습니다."})
		return
	}

	c.JSON(http.StatusOK, productResponses)
}

// CreateProduct godoc
// @Summary      상품 생성
// @Description  새로운 상품을 등록합니다. (관리자 전용)
// @Tags         products
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        productRequest body request.CreateProductRequest true "상품 정보"
// @Success      201 {object} response.ProductResponse "생성 성공"
// @Failure      400 {object} map[string]string "잘못된 요청"
// @Failure      403 {object} map[string]string "권한 없음"
// @Failure      500 {object} map[string]string "생성 실패"
// @Router       /products [post]
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

	responseData, err := pc.productInteractor.CreateProduct(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, responseData)
}

// UpdateStock godoc
// @Summary      재고 수정
// @Description  상품의 재고 수량을 수정합니다. (관리자 전용)
// @Tags         products
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        product_number path string true "상품 번호"
// @Param        stockRequest body request.UpdateStockRequest true "재고 정보"
// @Success      200 {object} map[string]string "수정 성공"
// @Failure      400 {object} map[string]string "잘못된 요청"
// @Failure      403 {object} map[string]string "권한 없음"
// @Failure      500 {object} map[string]string "수정 실패"
// @Router       /products/{product_number}/stock [put]
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

// DeleteProduct godoc
// @Summary      상품 삭제
// @Description  상품을 삭제합니다. (관리자 전용)
// @Tags         products
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        product_number path string true "상품 번호"
// @Success      200 {object} map[string]string "삭제 성공"
// @Failure      400 {object} map[string]string "잘못된 상품 번호"
// @Failure      403 {object} map[string]string "권한 없음"
// @Failure      500 {object} map[string]string "삭제 실패"
// @Router       /products/{product_number} [delete]
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
