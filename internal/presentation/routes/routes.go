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

	category := router.Group("/api/admin/categories")
	{
		category.POST("", categoryController.Create)
		category.GET("", categoryController.GetAll)
	}

	material := router.Group("/api/admin/materials")
	{
		material.POST("", materialController.Create)
		material.GET("", materialController.GetAll)
	}

	color := router.Group("/api/admin/colors")
	{
		color.POST("", colorController.Create)
		color.GET("", colorController.GetAll)
	}

	protected := router.Group("/api/protected")
	protected.Use(
		middleware.AuthMiddleware(jwtSecret),
	)
	{
		protected.GET("/test", func(c *gin.Context) {

			userID, _ := c.Get("user_id")
			role, _ := c.Get("role")

			c.JSON(200, gin.H{
				"user_id": userID,
				"role":    role,
			})
		})
	}
	admin := router.Group("/api/admin")
	admin.Use(
		middleware.AuthMiddleware(jwtSecret),
		middleware.RequireRole("admin"),
	)
	{
		admin.GET("/test", func(c *gin.Context) {

			c.JSON(200, gin.H{
				"message": "welcome admin",
			})

		})
	}

}
