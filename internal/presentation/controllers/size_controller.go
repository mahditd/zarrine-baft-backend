package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mahditd/zarrine-baft-backend/internal/application/services"
	"github.com/mahditd/zarrine-baft-backend/internal/presentation/dto"
)

type SizeController struct {
	sizeService *services.SizeService
}

func NewSizeController(
	sizeService *services.SizeService,
) *SizeController {

	return &SizeController{
		sizeService: sizeService,
	}
}

func (c *SizeController) GetAll(ctx *gin.Context) {

	sizes, err := c.sizeService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := make([]dto.SizeResponse, 0, len(sizes))

	for _, size := range sizes {
		response = append(response, dto.FromSize(&size))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"sizes": response,
	})
}
