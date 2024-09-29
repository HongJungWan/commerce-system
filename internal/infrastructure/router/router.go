package router

import (
	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/infrastructure/configs"
	controller "github.com/HongJungWan/commerce-system/internal/interfaces/controller"
	"github.com/HongJungWan/commerce-system/internal/usecases"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func NewRouter(conf configs.Config, db *gorm.DB) *gin.Engine {
	service := gin.New()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.AllowCredentials = true

	service.Use(cors.New(config))
	service.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "commerce system")
	})

	// 데이터베이스 마이그레이션, 일단 막아 두기
	db.Table("member").AutoMigrate(&domain.Member{})
	db.Table("order").AutoMigrate(&domain.Order{})
	db.Table("product").AutoMigrate(&domain.Product{})

	// Health Check 관련 설정
	healthCheckInteractor := usecases.NewHealthCheckInteractor()
	healthCheckController := controller.NewHealthCheckController(healthCheckInteractor)

	// FIXME: 비즈니스 로직 관련 설정하기

	router := service.Group("/api")

	// API 라우트 설정
	router.GET("/health", healthCheckController.HealthCheck)

	return service
}
