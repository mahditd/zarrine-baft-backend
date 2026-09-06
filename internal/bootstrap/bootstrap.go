package bootstrap

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/mahditd/zarrine-baft-backend/internal/application/services"
	"github.com/mahditd/zarrine-baft-backend/internal/config"
	"github.com/mahditd/zarrine-baft-backend/internal/infrastructure/database"
	"github.com/mahditd/zarrine-baft-backend/internal/infrastructure/repositories"
	"github.com/mahditd/zarrine-baft-backend/internal/infrastructure/storage"
	"github.com/mahditd/zarrine-baft-backend/internal/presentation/controllers"
	"github.com/mahditd/zarrine-baft-backend/internal/presentation/routes"
)

func Start() {

	cfg := config.Load()

	db := database.Connect(cfg)

	database.Migrate(db)

	database.SeedSizes(db)

	// User
	userRepository := repositories.NewUserRepository(db)

	authService := services.NewAuthService(
		userRepository,
		cfg.JWTSecret,
		cfg.JWTExpireHours,
	)

	authController := controllers.NewAuthController(
		authService,
	)

	// Category
	categoryRepository := repositories.NewCategoryRepository(db)

	categoryService := services.NewCategoryService(
		categoryRepository,
	)

	categoryController := controllers.NewCategoryController(
		categoryService,
	)

	// Material
	materialRepository := repositories.NewMaterialRepository(db)

	materialService := services.NewMaterialService(
		materialRepository,
	)

	materialController := controllers.NewMaterialController(
		materialService,
	)

	// Color
	colorRepository := repositories.NewColorRepository(db)

	colorService := services.NewColorService(
		colorRepository,
	)

	colorController := controllers.NewColorController(
		colorService,
	)

	// Size
	sizeRepository := repositories.NewSizeRepository(db)

	// Product
	productRepository := repositories.NewProductRepository(db)

	productService := services.NewProductService(
		productRepository,
		categoryRepository,
		materialRepository,
	)

	productController := controllers.NewProductController(
		productService,
	)

	// Product Variant
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

	localStorage := storage.NewLocalStorage(
		cfg.UploadPath + "/products",
	)

	// Product Image
	productImageRepository := repositories.NewProductImageRepository(db)

	productImageService := services.NewProductImageService(
		productImageRepository,
		productRepository,
		localStorage,
		cfg.BaseURL,
	)

	productImageController := controllers.NewProductImageController(
		productImageService,
	)

	// Product Request
	productRequestRepository := repositories.NewProductRequestRepository(db)

	productRequestService := services.NewProductRequestService(
		productRequestRepository,
		productVariantRepository,
		userRepository,
	)

	productRequestController := controllers.NewProductRequestController(
		productRequestService,
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
		productImageController,
		productRequestController,
		cfg.JWTSecret,
	)

	err := router.Run(
		fmt.Sprintf(":%s", cfg.AppPort),
	)

	if err != nil {
		panic(err)
	}
}
