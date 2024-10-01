package controller

import (
	"net/http"

	"github.com/HongJungWan/commerce-system/internal/usecases"
	"github.com/gin-gonic/gin"
)

type HealthCheckController struct {
	interactor usecases.HealthCheckInteractor
}

func NewHealthCheckController(interactor usecases.HealthCheckInteractor) *HealthCheckController {
	return &HealthCheckController{interactor}
}

// HealthCheck godoc
// @Summary      서비스 상태 확인
// @Description  서비스의 상태를 확인하고 정상 동작 여부를 검증합니다.
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200 {object} usecases.HealthStatus "서비스 상태"
// @Failure      500 {object} map[string]string "서버 오류"
// @Router       /health [get]
func (h *HealthCheckController) HealthCheck(c *gin.Context) {
	result := h.interactor.PerformHealthCheck()
	if result.Status != "Healthy" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "서비스에 문제가 발생했습니다."})
		return
	}
	c.JSON(http.StatusOK, result)
}
