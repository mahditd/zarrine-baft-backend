package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mahditd/zarrine-baft-backend/internal/application/services"
	"github.com/mahditd/zarrine-baft-backend/internal/presentation/dto"
)

type DashboardController struct {
	dashboardService *services.DashboardService
}

func NewDashboardController(
	dashboardService *services.DashboardService,
) *DashboardController {

	return &DashboardController{
		dashboardService: dashboardService,
	}
}

func (dc *DashboardController) GetDashboard(
	c *gin.Context,
) {

	dashboard, err := dc.dashboardService.GetDashboard()

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, dto.FromDashboardResult(dashboard))
}
