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

		//Slot
		InsertSlot (c echo.Context) error
		FindOneSlot (c echo.Context) error
		FindAllSlots (c echo.Context) error

		//Badminton
		InsertBadCourt ( c echo.Context) error
		FindCourt(c echo.Context) error
		InsertBadmintonSlot ( c echo.Context) error
		FindBadmintonSlot(c echo.Context) error
		
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


	res, err := h.facilityUsecase.FindManyFacility(ctx)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *facilityHttpHandler) InsertSlot(c echo.Context) error {
	facilityName := c.Param("facilityName")

	var slotRequest struct {
		StartTime       string `json:"start_time"`
		EndTime         string `json:"end_time"`
		MaxBookings     int    `json:"max_bookings"`
		CurrentBookings int    `json:"current_bookings"`
		FacilityType    string `json:"facility_type"`
	}

	if err := c.Bind(&slotRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: "+err.Error())
	}

	ctx := c.Request().Context()

	slot, err := h.facilityUsecase.InsertSlot(
		ctx,
		slotRequest.StartTime,
		slotRequest.EndTime,
		facilityName, 
		slotRequest.MaxBookings,
		slotRequest.CurrentBookings,
		slotRequest.FacilityType,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, slot)
}


func (h *facilityHttpHandler) FindOneSlot (c echo.Context) error {

	facilityName := c.QueryParam("facility_name") // Retrieve facilityName from query params

	if facilityName == "" {
		return response.ErrResponse(c, http.StatusBadRequest, "Facility name is required")
	}
	slotId := c.Param("slot_id")

	if slotId == "" {
		return response.ErrResponse(c, http.StatusBadRequest, "Slot id is required")
	}

	ctx := c.Request().Context()

	slot, err := h.facilityUsecase.FindOneSlot(ctx,facilityName, slotId)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, slot)
}

func (h *facilityHttpHandler) FindAllSlots(c echo.Context) error {
	ctx := c.Request().Context()
	facilityName := c.Param("facilityName")

	log.Printf("Finding slots for facility: %s", facilityName)

	slots, err := h.facilityUsecase.FindManySlot(ctx, facilityName)
	if err != nil {
		log.Printf("Error finding slots: %v", err)
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	log.Printf("Found %d slots for %s", len(slots), facilityName)

	// Return consistent response format
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data": slots,
	})
}

func (h *facilityHttpHandler) InsertBadCourt ( c echo.Context) error {
	ctx := c.Request().Context()

	var court facility.BadmintonCourt
	if err := c.Bind(&court); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	courtId, err := h.facilityUsecase.InsertBadCourt(ctx, &court)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, echo.Map{"court_id": courtId})
}

func (h *facilityHttpHandler) FindCourt(c echo.Context) error {
	ctx := c.Request().Context()

	courts, err := h.facilityUsecase.FindBadCourt(ctx)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data": courts,
	})
}

func (h *facilityHttpHandler) InsertBadmintonSlot ( c echo.Context) error {
    ctx := c.Request().Context()

	var slot facility.BadmintonSlot
	if err := c.Bind(&slot); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	slotId, err := h.facilityUsecase.InsertBadmintonSlot(ctx, &slot)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, echo.Map{"slot_id": slotId})
}

func (h *facilityHttpHandler) FindBadmintonSlot(c echo.Context) error {
	ctx := c.Request().Context()

	slots, err := h.facilityUsecase.FindBadmintonSlot(ctx)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data": slots,
	})
}