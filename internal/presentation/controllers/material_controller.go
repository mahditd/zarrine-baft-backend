package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mahditd/zarrine-baft-backend/internal/application/services"
	"github.com/mahditd/zarrine-baft-backend/internal/presentation/dto"
)

type MaterialController struct {
	materialService *services.MaterialService
}

func NewMaterialController(
	materialService *services.MaterialService,
) *MaterialController {

	return &MaterialController{
		materialService: materialService,
	}
}


type createMaterialRequest struct {
	NameFA string `json:"name_fa" binding:"required"`
	NameEN string `json:"name_en" binding:"required"`
}


func (c *MaterialController) Create(ctx *gin.Context) {

	var request createMaterialRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}


	material, err := c.materialService.Create(
		services.CreateMaterialInput{
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
		"message": "material created successfully",
		"material": dto.FromMaterial(material),
	})
}



func (c *MaterialController) GetAll(ctx *gin.Context) {

	materials, err := c.materialService.GetAll()

	if err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}


	response := make([]dto.MaterialResponse, 0)


	for _, material := range materials {

		response = append(
			response,
			dto.FromMaterial(&material),
		)
	}


	ctx.JSON(http.StatusOK, gin.H{
		"materials": response,
	})
}