package handler

import (
	"main/config"
	"main/modules/analytics/usecase"
	"main/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	AnalyticsHttpHandlerService interface {
		GetDashboardMetrics(c echo.Context) error
		GetUserAnalytics(c echo.Context) error
	}

	analyticsHttpHandler struct {
		cfg             *config.Config
		analyticsUsecase usecase.AnalyticsUsecaseService
	}
)

func NewAnalyticsHttpHandler(cfg *config.Config, analyticsUsecase usecase.AnalyticsUsecaseService) AnalyticsHttpHandlerService {
	return &analyticsHttpHandler{
		cfg:             cfg,
		analyticsUsecase: analyticsUsecase,
	}
}

func (h *analyticsHttpHandler) GetDashboardMetrics(c echo.Context) error {
	facilityName := c.QueryParam("facility_name")
	timeRange := c.QueryParam("time_range") // daily, weekly, monthly, yearly

	metrics, err := h.analyticsUsecase.GetDashboardMetrics(c.Request().Context(), facilityName, timeRange)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, metrics)
}

func (h *analyticsHttpHandler) GetUserAnalytics(c echo.Context) error {
	metrics, err := h.analyticsUsecase.GetUserAnalytics(c.Request().Context())
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, metrics)
} 