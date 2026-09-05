package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mahditd/zarrine-baft-backend/internal/application/services"
	"github.com/mahditd/zarrine-baft-backend/internal/presentation/dto"
)

type ColorController struct {
	colorService *services.ColorService
}

func NewColorController(
	colorService *services.ColorService,
) *ColorController {

	return &ColorController{
		colorService: colorService,
	}
}

type createColorRequest struct {
	NameFA  string `json:"name_fa" binding:"required"`
	NameEN  string `json:"name_en" binding:"required"`
	HexCode string `json:"hex_code" binding:"required"`
}

func (c *ColorController) Create(ctx *gin.Context) {

	var request createColorRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	color, err := c.colorService.Create(
		services.CreateColorInput{
			NameFA:  request.NameFA,
			NameEN:  request.NameEN,
			HexCode: request.HexCode,
		},
	)

	if err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "color created successfully",
		"color":   dto.FromColor(color),
	})
}

func (c *ColorController) GetAll(ctx *gin.Context) {

	colors, err := c.colorService.GetAll()

	if err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	response := make([]dto.ColorResponse, 0)

	for _, color := range colors {

		response = append(
			response,
			dto.FromColor(&color),
		)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"colors": response,
	})
}
