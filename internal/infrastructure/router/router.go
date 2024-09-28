package router

import (
	domain "github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/infrastructure/configs"
	controller "github.com/HongJungWan/commerce-system/internal/interfaces/controller"
	"github.com/HongJungWan/commerce-system/internal/usecases"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
		ctx.JSON(http.StatusOK, "ffmpeg video modules")
	})

	// 데이터베이스 마이그레이션
	db.Table("member").AutoMigrate(&domain.Member{})
	db.Table("order").AutoMigrate(&domain.Order{})
	db.Table("product").AutoMigrate(&domain.Product{})

	service.Static("/downloads", "./downloads")

	// Health Check 관련 설정
	healthCheckInteractor := usecases.NewHealthCheckInteractor()
	healthCheckController := controller.NewHealthCheckController(healthCheckInteractor)

	// FIXME: 비즈니스 로직 관련 설정하기

	router := service.Group("/api")

	// API 라우트 설정
	router.GET("/health", healthCheckController.HealthCheck)

	// Swagger UI 라우트 추가
	service.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return service
}
