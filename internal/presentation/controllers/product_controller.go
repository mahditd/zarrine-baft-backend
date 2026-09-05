package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/mahditd/zarrine-baft-backend/internal/application/services"
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

	products, err := pc.productService.GetAll()

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	response := make([]dto.ProductResponse, 0)

	for _, product := range products {

		response = append(
			response,
			dto.FromProduct(&product),
		)

	}

	c.JSON(http.StatusOK, gin.H{
		"products": response,
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

	err = pc.productService.Delete(uint(id))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "product deleted successfully",
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

	categoryID, _ := strconv.ParseUint(
		c.Query("category_id"),
		10,
		64,
	)

	materialID, _ := strconv.ParseUint(
		c.Query("material_id"),
		10,
		64,
	)

	colorID, _ := strconv.ParseUint(
		c.Query("color_id"),
		10,
		64,
	)

	products, total, err := pc.productService.GetActiveProducts(
		page,
		limit,
		uint(categoryID),
		uint(materialID),
		uint(colorID),
	)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	response := make([]dto.ProductResponse, 0)

	for _, product := range products {

		response = append(
			response,
			dto.FromProduct(&product),
		)

	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, gin.H{

		"products": response,

		"page": page,

		"limit": limit,

		"total_products": total,

		"total_pages": totalPages,
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
