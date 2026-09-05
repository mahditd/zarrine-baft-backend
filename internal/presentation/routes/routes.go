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
	jwtSecret string,
) {

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	auth := router.Group("/api/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	admin := router.Group("/api/admin")

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
		admin.DELETE("/product-images/:id", productImageController.Delete)

		admin.POST("/products", productController.Create)
		admin.GET("/products", productController.GetAll)

		admin.POST("/product-variants", productVariantController.Create)
		admin.GET("/products/:id/variants", productVariantController.GetByProductID)

		admin.GET("/test", func(c *gin.Context) {

			c.JSON(200, gin.H{
				"message": "welcome admin",
			})

		})
	}

}
