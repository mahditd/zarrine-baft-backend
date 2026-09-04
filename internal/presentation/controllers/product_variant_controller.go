package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mahditd/zarrine-baft-backend/internal/application/services"
	"github.com/mahditd/zarrine-baft-backend/internal/presentation/dto"
)


type ProductVariantController struct {
	productVariantService *services.ProductVariantService
}


func NewProductVariantController(
	productVariantService *services.ProductVariantService,
) *ProductVariantController {

	return &ProductVariantController{
		productVariantService: productVariantService,
	}
}



func (pvc *ProductVariantController) Create(
	c *gin.Context,
) {

	var input services.CreateProductVariantInput


	if err := c.ShouldBindJSON(&input); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}



	variant, err := pvc.productVariantService.Create(input)


	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}



	c.JSON(http.StatusCreated, gin.H{

		"message": "product variant created successfully",

		"variant": dto.FromProductVariant(variant),

	})
}



func (pvc *ProductVariantController) GetByProductID(
	c *gin.Context,
) {

	productID := c.Param("id")


	variantList, err := pvc.productVariantService.GetByProductID(
		parseUint(productID),
	)


	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}



	response := make([]dto.ProductVariantResponse, 0)


	for _, variant := range variantList {

		response = append(
			response,
			dto.FromProductVariant(&variant),
		)

	}



	c.JSON(http.StatusOK, gin.H{
		"variants": response,
	})

}