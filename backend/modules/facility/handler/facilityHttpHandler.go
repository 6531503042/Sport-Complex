package handler

import (
	"log"
	"main/config"
	"main/modules/facility"
	"main/modules/facility/usecase"
	"main/pkg/request"
	"main/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	NewFacilityHttpHandlerService interface {
		CreateFacility(c echo.Context) error
		FindOneFacility(c echo.Context) error
		FindManyFacility(c echo.Context) error
	}

	facilityHttpHandler struct {
		cfg            *config.Config
		facilityUsecase usecase.FacilityUsecaseService
	}
)

func NewFacilityHttpHandler(cfg *config.Config, facilityUsecase usecase.FacilityUsecaseService) NewFacilityHttpHandlerService {
	return &facilityHttpHandler{cfg: cfg, facilityUsecase: facilityUsecase}
}

func (h *facilityHttpHandler) CreateFacility(c echo.Context) error {

	log.Println("Received request to create facility")

	ctx := c.Request().Context()

	// Use the custom binding
	wrapper := request.ContextWrapper(c)

	req := new(facility.CreateFaciliityRequest) // Bind the request body to the struct

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "Invalid request payload")
	}

	// Pass the facility name for repository operations
	res, err := h.facilityUsecase.CreateFacility(ctx, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, res)
}

func (h *facilityHttpHandler) FindOneFacility(c echo.Context) error {

	log.Println("Received request to find one facility")

	ctx := c.Request().Context()

	facilityId := c.Param("facility_id")
	facilityName := c.QueryParam("facility_name") // Retrieve facilityName from query params

	if facilityName == "" {
		return response.ErrResponse(c, http.StatusBadRequest, "Facility name is required")
	}

	res, err := h.facilityUsecase.FindOneFacility(ctx, facilityId, facilityName)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *facilityHttpHandler) FindManyFacility(c echo.Context) error {

	log.Println("Received request to find many facilities")

	ctx := c.Request().Context()

	facilityName := c.QueryParam("facility_name") // Retrieve facilityName from query params

	if facilityName == "" {
		return response.ErrResponse(c, http.StatusBadRequest, "Facility name is required")
	}

	res, err := h.facilityUsecase.FindManyFacility(ctx, facilityName)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}
