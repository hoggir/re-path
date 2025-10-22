package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoggir/re-path/redirect-service/internal/dto"
	"github.com/hoggir/re-path/redirect-service/internal/service"
)

type DashboardHandler struct {
	dashboardService service.DashboardService
}

func NewDashboardHandler(dashboardService service.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// Dashboard endpoint
// @Summary Get dashboard statistics
// @Description Get dashboard statistics (requires authentication)
// @Tags Dashboard
// @Security BearerAuth
// @Produce json
// @Accept json
// @Success 200 {object} dto.Response{data=dto.DashboardResponse}
// @Failure 401 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Router /api/dashboard [get]
func (d *DashboardHandler) GetDashboardByShortUrl(c *gin.Context) {
	userID, _ := c.Get("user_id")

	dashboardData, err := d.dashboardService.GetDashboard(c.Request.Context(), userID.(int))
	if err != nil {
		fmt.Printf("❌ Error getting dashboard: %v\n", err)
		dto.HandleError(c, err)
		return
	}

	response := dto.DashboardResponse{
		TotalLink:    dashboardData.TotalLinks,
		TotalClick:   dashboardData.TotalClicks,
		UniqVisitors: dashboardData.UniqVisitors,
		TopLinks:     dashboardData.TopLinks,
		StatLinks:    dashboardData.StatLinks,
	}

	if dashboardData.IsLimited() {
		fmt.Printf("⚠️ Dashboard data is limited: %s\n", dashboardData.GetMessage())
	}

	dto.SuccessResponse(c, http.StatusOK, "Dashboard retrieved successfully", response)
}
