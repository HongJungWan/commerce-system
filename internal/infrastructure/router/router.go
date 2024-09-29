package router

import (
	"github.com/HongJungWan/commerce-system/internal/interfaces/middleware"
	"net/http"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/infrastructure/configs"
	"github.com/HongJungWan/commerce-system/internal/infrastructure/repository"
	"github.com/HongJungWan/commerce-system/internal/interfaces/controller"
	"github.com/HongJungWan/commerce-system/internal/usecases"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	// 데이터베이스 마이그레이션
	db.AutoMigrate(&domain.Member{})
	db.AutoMigrate(&domain.Product{})
	db.AutoMigrate(&domain.Order{})

	// Health Check 관련 설정
	healthCheckInteractor := usecases.NewHealthCheckInteractor()
	healthCheckController := controller.NewHealthCheckController(healthCheckInteractor)

	// 회원 관련 설정
	memberRepo := repository.NewMemberRepository(db)
	memberInteractor := usecases.NewMemberInteractor(memberRepo)
	authInteractor := usecases.NewAuthUseCase("commerce-system", memberRepo)
	memberController := controller.NewMemberController(memberInteractor, authInteractor)
	authController := controller.NewAuthController(authInteractor)

	// JWT 미들웨어 설정
	authMiddleware := middleware.JWTAuthMiddleware()

	router := service.Group("/api")

	// 헬스체크 엔드포인트 설정
	router.GET("/health", healthCheckController.HealthCheck)

	// Auth 엔드포인트 설정
	router.POST("/login", authController.Login)

	// 회원 엔드포인트 설정
	router.POST("/members", memberController.Register)
	router.GET("/members/me", authMiddleware, memberController.GetMyInfo)
	router.PUT("/members/me", authMiddleware, memberController.UpdateMyInfo)
	router.DELETE("/members/me", authMiddleware, memberController.DeleteMyAccount)
	router.GET("/members", authMiddleware, memberController.GetAllMembers)
	router.GET("/members/stats", authMiddleware, memberController.GetMemberStats)

	return service
}
