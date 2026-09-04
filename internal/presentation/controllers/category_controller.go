package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mahditd/zarrine-baft-backend/internal/application/services"
	"github.com/mahditd/zarrine-baft-backend/internal/presentation/dto"
)

type CategoryController struct {
	categoryService *services.CategoryService
}

func NewCategoryController(
	categoryService *services.CategoryService,
) *CategoryController {

	return &CategoryController{
		categoryService: categoryService,
	}
}

type createCategoryRequest struct {
	NameFA string `json:"name_fa" binding:"required"`
	NameEN string `json:"name_en" binding:"required"`
}

func (c *CategoryController) Create(ctx *gin.Context) {

	var request createCategoryRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	category, err := c.categoryService.Create(
		services.CreateCategoryInput{
			NameFA: request.NameFA,
			NameEN: request.NameEN,
		},
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":  "category created successfully",
		"category": dto.FromCategory(category),
	})
}

func (c *CategoryController) GetAll(ctx *gin.Context) {

	categories, err := c.categoryService.GetAll()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}


	response := make([]dto.CategoryResponse, 0)

	for _, category := range categories {

		response = append(
			response,
			dto.FromCategory(&category),
		)
	}


	ctx.JSON(http.StatusOK, gin.H{
		"categories": response,
	})
}