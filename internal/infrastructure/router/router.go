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
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	service.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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
	authInteractor := usecases.NewAuthUseCase("commerce-system", memberRepo)
	memberInteractor := usecases.NewMemberInteractor(memberRepo, authInteractor)
	memberController := controller.NewMemberController(memberInteractor, authInteractor)
	authController := controller.NewAuthController(authInteractor)

	// 상품 관련 설정
	productRepo := repository.NewProductRepository(db)
	productInteractor := usecases.NewProductInteractor(productRepo, db)
	productController := controller.NewProductController(productInteractor)

	// 주문 관련 설정
	orderRepo := repository.NewOrderRepository(db)
	orderInteractor := usecases.NewOrderInteractor(orderRepo, memberRepo, productRepo)
	orderController := controller.NewOrderController(orderInteractor)

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

	// 상품 엔드포인트 설정
	router.GET("/products", productController.GetProducts)
	router.POST("/products", authMiddleware, productController.CreateProduct)
	router.PUT("/products/:id/stock", authMiddleware, productController.UpdateStock)
	router.DELETE("/products/:id", authMiddleware, productController.DeleteProduct)

	// 주문 엔드포인트 설정
	router.POST("/orders", authMiddleware, orderController.CreateOrder)
	router.GET("/orders/me", authMiddleware, orderController.GetMyOrders)
	router.PUT("/orders/:id/cancel", authMiddleware, orderController.CancelOrder)
	router.GET("/orders/stats", authMiddleware, orderController.GetMonthlyStats)

	return service
}
