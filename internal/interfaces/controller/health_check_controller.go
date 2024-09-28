package controller

import (
	"github.com/HongJungWan/commerce-system/internal/usecases"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthCheckController struct {
	interactor usecases.HealthCheckInteractor
}

func NewHealthCheckController(interactor usecases.HealthCheckInteractor) *HealthCheckController {
	return &HealthCheckController{interactor}
}

// HealthCheck godoc
// @Summary      Performs a health check
// @Description  Endpoint to perform a health check and verify the service status
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200 {object} usecases.HealthStatus "Health check result"
// @Router       /health [get]
func (h *HealthCheckController) HealthCheck(c *gin.Context) {
	result := h.interactor.PerformHealthCheck()
	c.JSON(http.StatusOK, result)
}
