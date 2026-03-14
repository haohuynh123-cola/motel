package handler

import (
	"net/http"
	"tro-go/internal/port"

	"github.com/labstack/echo/v4"
)

type DashboardHandler struct {
	dashboardUC port.DashboardUseCase
}

func NewDashboardHandler(g *echo.Group, dashboardUC port.DashboardUseCase) {
	h := &DashboardHandler{dashboardUC: dashboardUC}

	g.GET("/dashboard/stats", h.GetStats)
}

func (h *DashboardHandler) GetStats(c echo.Context) error {
	response, err := h.dashboardUC.GetDashboardStats(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, port.ApiResponse{Status: false, Data: err.Error()})
	}
	return c.JSON(http.StatusOK, response)
}
