package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/mahditd/zarrine-baft-backend/internal/application/services"
	"github.com/mahditd/zarrine-baft-backend/internal/domain/repositories"
	"github.com/mahditd/zarrine-baft-backend/internal/presentation/dto"
)

type ProductController struct {
	productService *services.ProductService
}

func NewProductController(
	productService *services.ProductService,
) *ProductController {

	return &ProductController{
		productService: productService,
	}
}

func parseUintSlice(c *gin.Context, key string) []uint {
	values := c.QueryArray(key)
	single := c.Query(key)
	if single != "" && len(values) <= 1 {
		parts := strings.Split(single, ",")
		values = parts
	}

	result := make([]uint, 0)
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if id, err := strconv.ParseUint(v, 10, 64); err == nil {
			result = append(result, uint(id))
		}
	}
	return result
}

func (pc *ProductController) Create(
	c *gin.Context,
) {

	var input services.CreateProductInput

	if err := c.ShouldBindJSON(&input); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	product, err := pc.productService.Create(input)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{

		"message": "product created successfully",

		"product": dto.FromProduct(product),
	})
}

func (pc *ProductController) GetAll(
	c *gin.Context,
) {

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit := 24
	search := strings.TrimSpace(c.Query("search"))

	var isActive *bool
	if activeStr := strings.TrimSpace(c.Query("is_active")); activeStr != "" {
		val := activeStr == "true" || activeStr == "1"
		isActive = &val
	}

	filter := repositories.ProductFilter{
		Page:     page,
		Limit:    limit,
		Search:   search,
		IsActive: isActive,
	}

	products, total, err := pc.productService.GetAdminProducts(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := make([]dto.ProductResponse, 0, len(products))
	for _, product := range products {
		response = append(response, dto.FromProduct(&product))
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, gin.H{
		"products":       response,
		"page":           page,
		"limit":          limit,
		"total_products": total,
		"total_pages":    totalPages,
	})
}

func (pc *ProductController) GetByID(
	c *gin.Context,
) {

	id, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid product id",
			},
		)

		return
	}

	product, err := pc.productService.GetByID(
		uint(id),
	)

	if err != nil {

		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"product": dto.FromProduct(product),
		},
	)
}

func (pc *ProductController) Update(
	c *gin.Context,
) {

	id, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid product id",
			},
		)

		return
	}

	var input services.UpdateProductInput

	if err := c.ShouldBindJSON(&input); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	product, err := pc.productService.Update(
		uint(id),
		input,
	)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "product updated successfully",
			"product": dto.FromProduct(product),
		},
	)
}

func (pc *ProductController) UpdateStatus(
	c *gin.Context,
) {

	id, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product id",
		})
		return
	}

	var input struct {
		IsActive bool `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = pc.productService.UpdateStatus(
		uint(id),
		input.IsActive,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "product status updated successfully",
	})
}

func (pc *ProductController) Delete(
	c *gin.Context,
) {
	// SRS Section 6: Products cannot be deleted.
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "products cannot be deleted; use deactivate instead",
	})
}

func (pc *ProductController) Reorder(
	c *gin.Context,
) {
	var input struct {
		ProductIDs []uint `json:"product_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := pc.productService.Reorder(input.ProductIDs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "products reordered successfully",
	})
}

func (pc *ProductController) GetActiveProducts(
	c *gin.Context,
) {

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid page number",
		})
		return
	}

	limit := 24
	search := strings.TrimSpace(c.Query("search"))

	categoryIDs := parseUintSlice(c, "category_ids")
	if len(categoryIDs) == 0 {
		categoryIDs = parseUintSlice(c, "category_id")
	}

	materialIDs := parseUintSlice(c, "material_ids")
	if len(materialIDs) == 0 {
		materialIDs = parseUintSlice(c, "material_id")
	}

	colorIDs := parseUintSlice(c, "color_ids")
	if len(colorIDs) == 0 {
		colorIDs = parseUintSlice(c, "color_id")
	}

	sizeIDs := parseUintSlice(c, "size_ids")
	if len(sizeIDs) == 0 {
		sizeIDs = parseUintSlice(c, "size_id")
	}

	filter := repositories.ProductFilter{
		Page:        page,
		Limit:       limit,
		Search:      search,
		CategoryIDs: categoryIDs,
		MaterialIDs: materialIDs,
		ColorIDs:    colorIDs,
		SizeIDs:     sizeIDs,
	}

	products, total, err := pc.productService.GetActiveProducts(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := make([]dto.ProductResponse, 0, len(products))
	for _, product := range products {
		response = append(response, dto.FromProduct(&product))
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, gin.H{
		"products":       response,
		"page":           page,
		"limit":          limit,
		"total_products": total,
		"total_pages":    totalPages,
	})
}

func (pc *ProductController) GetActiveByID(
	c *gin.Context,
) {

	id, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product id",
		})
		return
	}

	product, err := pc.productService.GetActiveByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "product not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": dto.FromProduct(product),
	})
}

