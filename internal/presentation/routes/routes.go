package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/mahditd/zarrine-baft-backend/internal/presentation/controllers"
	"github.com/mahditd/zarrine-baft-backend/internal/presentation/middleware"
)

func SetupRoutes(
	router *gin.Engine,
	authController *controllers.AuthController,
	categoryController *controllers.CategoryController,
	materialController *controllers.MaterialController,
	colorController *controllers.ColorController,
	productController *controllers.ProductController,
	productVariantController *controllers.ProductVariantController,
	productImageController *controllers.ProductImageController,
	productRequestController *controllers.ProductRequestController,
	dashboardController *controllers.DashboardController,
	jwtSecret string,
	authLimiter *middleware.IPRateLimiter,
	apiLimiter *middleware.IPRateLimiter,
) {

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	router.Static("/uploads", "./uploads")

	auth := router.Group("/api/auth")
	if authLimiter != nil {
		auth.Use(middleware.RateLimitMiddleware(authLimiter))
	}
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	// Authenticated Customer / User routes
	authenticated := router.Group("/api")
	if apiLimiter != nil {
		authenticated.Use(middleware.RateLimitMiddleware(apiLimiter))
	}
	authenticated.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// Profile (SRS 3.3)
		authenticated.GET("/me/profile", authController.GetProfile)
		authenticated.PATCH("/me/profile", authController.UpdateProfile)

		// Requests (SRS 12 & 14)
		authenticated.POST("/requests", productRequestController.Create)
		authenticated.GET("/me/requests", productRequestController.GetMyRequests)
		authenticated.GET("/me/requests/:id", productRequestController.GetMyRequestByID)
		authenticated.PATCH("/me/requests/:id/cancel", productRequestController.Cancel)
	}

	// Public catalog routes with general API rate limiting (SRS 22)
	if apiLimiter != nil {
		router.GET("/api/products", middleware.RateLimitMiddleware(apiLimiter), productController.GetActiveProducts)
		router.GET("/api/products/:id", middleware.RateLimitMiddleware(apiLimiter), productController.GetActiveByID)
	} else {
		router.GET("/api/products", productController.GetActiveProducts)
		router.GET("/api/products/:id", productController.GetActiveByID)
	}

	admin := router.Group("/api/admin")

	if apiLimiter != nil {
		admin.Use(middleware.RateLimitMiddleware(apiLimiter))
	}
	admin.Use(
		middleware.AuthMiddleware(jwtSecret),
		middleware.RequireRole("admin"),
	)

	{
		admin.POST("/categories", categoryController.Create)
		admin.GET("/categories", categoryController.GetAll)

		admin.POST("/materials", materialController.Create)
		admin.GET("/materials", materialController.GetAll)

		admin.POST("/colors", colorController.Create)
		admin.GET("/colors", colorController.GetAll)

		admin.POST("/products/:id/images", productImageController.Create)
		admin.GET("/products/:id/images", productImageController.GetByProductID)
		admin.PUT("/products/:id/images/:imageId", productImageController.Replace)
		admin.DELETE("/product-images/:id", productImageController.Delete)
		admin.PATCH("/products/:id/images/reorder", productImageController.Reorder)

		admin.POST("/products", productController.Create)
		admin.GET("/products", productController.GetAll)
		admin.PATCH("/products/reorder", productController.Reorder)
		admin.GET("/products/:id", productController.GetByID)
		admin.PATCH("/products/:id", productController.Update)
		admin.PATCH("/products/:id/status", productController.UpdateStatus)
		admin.DELETE("/products/:id", productController.Delete)

		admin.POST("/product-variants", productVariantController.Create)
		admin.PATCH("/product-variants/:id", productVariantController.Update)
		admin.DELETE("/product-variants/:id", productVariantController.Delete)
		admin.GET("/products/:id/variants", productVariantController.GetByProductID)

		admin.GET(
			"/dashboard",
			dashboardController.GetDashboard,
		)

		// Product Requests
		admin.GET(
			"/requests",
			productRequestController.GetAll,
		)
		admin.GET(
			"/requests/:id",
			productRequestController.GetByID,
		)

		admin.PATCH(
			"/requests/:id/status",
			productRequestController.UpdateStatus,
		)

		admin.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "welcome admin",
			})
		})
	}

}
