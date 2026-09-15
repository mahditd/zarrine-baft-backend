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
	"github.com/mahditd/zarrine-baft-backend/internal/presentation/middleware"
	"github.com/mahditd/zarrine-baft-backend/internal/presentation/routes"
)

func Start() {

	cfg := config.Load()

	db := database.Connect(cfg)

	database.Migrate(db)

	database.SeedSizes(db)

	if err := database.SeedAdmin(db, cfg.AdminPhone, cfg.AdminPassword, cfg.AdminName); err != nil {
		fmt.Printf("Warning: failed to seed admin user: %v\n", err)
	}

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

	// Product Image repository is needed by ProductService (min-1 image on activate).
	productImageRepository := repositories.NewProductImageRepository(db)

	productService := services.NewProductService(
		productRepository,
		categoryRepository,
		materialRepository,
		productImageRepository,
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

	// Product Image (repository already created above for ProductService)
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

	// Dashboard
	dashboardService := services.NewDashboardService(
		productRepository,
		productRequestRepository,
	)

	dashboardController := controllers.NewDashboardController(
		dashboardService,
	)

	router := gin.Default()

	// Baseline security headers (SRS 22). HTTPS is terminated at the
	// production reverse proxy in front of this service.
	router.Use(middleware.SecurityHeaders())

	// CORS support for web frontend SPA (SRS Technical Architecture)
	router.Use(middleware.CORSMiddleware(cfg.ClientURL))

	// Rate limiters (SRS 22: Login protection & API protection)
	authLimiter := middleware.NewIPRateLimiter(15, 5)   // 15 auth requests/min with burst 5
	apiLimiter := middleware.NewIPRateLimiter(120, 30)  // 120 API requests/min with burst 30

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
		dashboardController,
		cfg.JWTSecret,
		authLimiter,
		apiLimiter,
	)

	err := router.Run(
		fmt.Sprintf(":%s", cfg.AppPort),
	)

	if err != nil {
		panic(err)
	}
}
