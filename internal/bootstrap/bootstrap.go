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

	database.SeedSizes(db)

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

	materialRepository := repositories.NewMaterialRepository(db)

	materialService := services.NewMaterialService(
		materialRepository,
	)

	materialController := controllers.NewMaterialController(
		materialService,
	)

	colorRepository := repositories.NewColorRepository(db)

	colorService := services.NewColorService(
		colorRepository,
	)

	colorController := controllers.NewColorController(
		colorService,
	)

	sizeRepository := repositories.NewSizeRepository(db)

	productRepository := repositories.NewProductRepository(db)

	productService := services.NewProductService(
		productRepository,
		categoryRepository,
		materialRepository,
	)

	productController := controllers.NewProductController(
		productService,
	)
	productVariantRepository := repositories.NewProductVariantRepository(db)

	productVariantService := services.NewProductVariantService(
		productVariantRepository,
		productRepository,
		colorRepository,
		sizeRepository,
	)

	productVariantController := controllers.NewProductVariantController(
		productVariantService,
	)

	router := gin.Default()

	routes.SetupRoutes(
		router,
		authController,
		categoryController,
		materialController,
		colorController,
		productController,
		productVariantController,
		cfg.JWTSecret,
	)

	router.Run(
		fmt.Sprintf(":%s", cfg.AppPort),
	)
}
