package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mahditd/zarrine-baft-backend/internal/presentation/dto"

	"github.com/mahditd/zarrine-baft-backend/internal/application/services"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(
	authService *services.AuthService,
) *AuthController {

	return &AuthController{
		authService: authService,
	}
}

type registerRequest struct {
	FullName        string `json:"full_name" binding:"required"`
	Phone           string `json:"phone" binding:"required"`
	Email           string `json:"email"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	CompanyName     string `json:"company_name"`
	CompanyPhone    string `json:"company_phone"`
	Country         string `json:"country"`
	Address         string `json:"address"`
}

func (c *AuthController) Register(ctx *gin.Context) {

	var request registerRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := c.authService.Register(
		services.RegisterInput{
			FullName:        request.FullName,
			Phone:           request.Phone,
			Email:           request.Email,
			Password:        request.Password,
			ConfirmPassword: request.ConfirmPassword,
			CompanyName:     request.CompanyName,
			CompanyPhone:    request.CompanyPhone,
			Country:         request.Country,
			Address:         request.Address,
		},
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
		"user":    dto.FromUser(user),
	})
}

func (c *AuthController) Login(ctx *gin.Context) {

	var request dto.LoginRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := c.authService.Login(
		services.LoginInput{
			Phone:    request.Phone,
			Password: request.Password,
		},
	)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   result.Token,
		"user":    dto.FromUser(result.User),
	})
}

func (c *AuthController) GetProfile(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	userID := userIDVal.(uint)

	user, err := c.authService.GetProfile(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": dto.FromUser(user),
	})
}

type updateProfileRequest struct {
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	CompanyName  string `json:"company_name"`
	CompanyPhone string `json:"company_phone"`
	Country      string `json:"country"`
	Address      string `json:"address"`
}

func (c *AuthController) UpdateProfile(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	userID := userIDVal.(uint)

	var request updateProfileRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := c.authService.UpdateProfile(
		userID,
		services.UpdateProfileInput{
			FullName:     request.FullName,
			Email:        request.Email,
			CompanyName:  request.CompanyName,
			CompanyPhone: request.CompanyPhone,
			Country:      request.Country,
			Address:      request.Address,
		},
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "profile updated successfully",
		"user":    dto.FromUser(user),
	})
}

