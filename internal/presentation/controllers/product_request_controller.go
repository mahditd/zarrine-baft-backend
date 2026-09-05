package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mahditd/zarrine-baft-backend/internal/application/services"
	"github.com/mahditd/zarrine-baft-backend/internal/presentation/dto"
)

type ProductRequestController struct {
	service *services.ProductRequestService
}

func NewProductRequestController(
	service *services.ProductRequestService,
) *ProductRequestController {

	return &ProductRequestController{
		service: service,
	}
}

func (c *ProductRequestController) Create(
	ctx *gin.Context,
) {

	var input services.CreateProductRequestInput

	if err := ctx.ShouldBindJSON(&input); err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	request, err := c.service.Create(input)

	if err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "product request created successfully",
		"request": dto.FromProductRequest(request),
	})
}

func (c *ProductRequestController) GetAll(
	ctx *gin.Context,
) {

	requests, err := c.service.GetAll()

	if err != nil {

		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"requests": dto.FromProductRequests(requests),
		},
	)
}
