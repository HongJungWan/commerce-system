package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HongJungWan/commerce-system/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckController_Integration(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)

	healthCheckInteractor := usecases.NewHealthCheckInteractor()
	healthCheckController := NewHealthCheckController(healthCheckInteractor)

	router := gin.Default()
	router.GET("/api/health", healthCheckController.HealthCheck)

	// When
	req, err := http.NewRequest(http.MethodGet, "/api/health", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Then
	assert.Equal(t, http.StatusOK, rr.Code, "/api/health")
	assert.JSONEq(t, `{"status": "Healthy", "message": "Success"}`, rr.Body.String())
}
