package controllers

import (
	"net/http"

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