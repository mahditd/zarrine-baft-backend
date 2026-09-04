package bootstrap

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/mahditd/zarrine-baft-backend/internal/application/services"
	"github.com/mahditd/zarrine-baft-backend/internal/config"
	"github.com/mahditd/zarrine-baft-backend/internal/infrastructure/database"
	"github.com/mahditd/zarrine-baft-backend/internal/infrastructure/repositories"
	"github.com/mahditd/zarrine-baft-backend/internal/presentation/controllers"
	"github.com/mahditd/zarrine-baft-backend/internal/presentation/routes"
)

func Start() {

	cfg := config.Load()

	db := database.Connect(cfg)

	fmt.Println(db)

	database.Migrate(db)

	userRepository := repositories.NewUserRepository(db)

	authService := services.NewAuthService(
		userRepository,
		cfg.JWTSecret,
		cfg.JWTExpireHours,
	)

	authController := controllers.NewAuthController(
		authService,
	)

	categoryRepository := repositories.NewCategoryRepository(db)

	categoryService := services.NewCategoryService(
		categoryRepository,
	)

	categoryController := controllers.NewCategoryController(
		categoryService,
	)

	router := gin.Default()

	routes.SetupRoutes(
		router,
		authController,
		categoryController,
		cfg.JWTSecret,
	)

	router.Run(
		fmt.Sprintf(":%s", cfg.AppPort),
	)
}
