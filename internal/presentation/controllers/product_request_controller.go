package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mahditd/zarrine-baft-backend/internal/application/services"
	"github.com/mahditd/zarrine-baft-backend/internal/presentation/dto"
	"strings"
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

func (c *ProductRequestController) GetByID(
	ctx *gin.Context,
) {

	id, err := parseUint(ctx.Param("id"))

	if err != nil {

		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid request id",
			},
		)

		return
	}

	request, err := c.service.GetByID(id)

	if err != nil {

		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"request": dto.FromProductRequest(request),
		},
	)
}

func (c *ProductRequestController) UpdateStatus(
	ctx *gin.Context,
) {

	id, err := parseUint(ctx.Param("id"))

	if err != nil {

		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid request id",
			},
		)

		return
	}

	var input struct {
		Status string `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {

		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	input.Status = strings.TrimSpace(input.Status)

	err = c.service.UpdateStatus(
		id,
		input.Status,
	)

	if err != nil {

		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "request status updated successfully",
		},
	)
}
