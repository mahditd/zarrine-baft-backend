package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/mahditd/zarrine-baft-backend/internal/application/services"
	"github.com/mahditd/zarrine-baft-backend/internal/presentation/dto"
)

type ProductImageController struct {
	service *services.ProductImageService
}

func NewProductImageController(
	service *services.ProductImageService,
) *ProductImageController {

	return &ProductImageController{
		service: service,
	}
}

func (c *ProductImageController) Create(
	ctx *gin.Context,
) {

	productID, err := strconv.ParseUint(
		ctx.Param("id"),
		10,
		64,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product id",
		})
		return
	}

	file, err := ctx.FormFile("image")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "image file is required",
		})
		return
	}

	image, err := c.service.Upload(
		services.UploadProductImageInput{
			ProductID: uint(productID),
			File:      file,
		},
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "product image uploaded successfully",
		"image":   dto.FromProductImage(image),
	})
}

func (c *ProductImageController) GetByProductID(
	ctx *gin.Context,
) {

	productID, err := strconv.ParseUint(
		ctx.Param("id"),
		10,
		64,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product id",
		})
		return
	}

	images, err := c.service.GetByProductID(
		uint(productID),
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := make([]dto.ProductImageResponse, 0)

	for _, image := range images {
		response = append(
			response,
			dto.FromProductImage(&image),
		)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"images": response,
	})
}

func (c *ProductImageController) Delete(
	ctx *gin.Context,
) {

	id, err := strconv.ParseUint(
		ctx.Param("id"),
		10,
		64,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid image id",
		})
		return
	}

	err = c.service.Delete(
		uint(id),
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "product image deleted successfully",
	})
}

func (c *ProductImageController) Reorder(
	ctx *gin.Context,
) {

	productID, err := strconv.ParseUint(
		ctx.Param("id"),
		10,
		64,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product id",
		})
		return
	}

	var body struct {
		ImageIDs []uint `json:"image_ids"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = c.service.Reorder(
		uint(productID),
		body.ImageIDs,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "product images reordered successfully",
	})
}
