package handler

import (
	"log"
	"main/config"
	"main/modules/booking"
	"main/modules/booking/usecase"
	"main/pkg/request"
	"main/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type(
	NewSlotHttpHandlerService interface {
		InsertSlot(c echo.Context) error
		UpdateSlot(c echo.Context) error
		FindSlot(c echo.Context) error
		FindAllSlots(c echo.Context) error
		EnableOrDisableSlot(c echo.Context) error
	}

	slotsHttpHandler struct {
		cfg *config.Config
		slotsUsecase usecase.SlotUsecaseService
		
	}
)

func NewSlotHttpHandler(cfg *config.Config, slotUsecase usecase.SlotUsecaseService) NewSlotHttpHandlerService {
	return &slotsHttpHandler{cfg: cfg, slotsUsecase: slotUsecase}
}


func (h *slotsHttpHandler) InsertSlot(c echo.Context) error {
    log.Println("Received request to create slot")

    ctx := c.Request().Context()
    wrapper := request.ContextWrapper(c)

    req := new(booking.Slot)
    if err := wrapper.Bind(req); err != nil {
        return response.ErrResponse(c, http.StatusBadRequest, "Invalid request payload")
    }

    // Since req.StartTime and req.EndTime are now strings, no need to convert them to time.Time
    startTime := req.StartTime // string in "HH:mm"
    endTime := req.EndTime     // string in "HH:mm"

    // Validate that startTime and endTime are not empty
    if startTime == "" || endTime == "" {
        return response.ErrResponse(c, http.StatusBadRequest, "Start time or end time cannot be empty")
    }

    // Pass the time strings to the usecase
    slot, err := h.slotsUsecase.InsertSlot(ctx, startTime, endTime)
    if err != nil {
        return response.ErrResponse(c, http.StatusBadRequest, err.Error())
    }

    return response.SuccessResponse(c, http.StatusCreated, slot)
}




func (h *slotsHttpHandler) UpdateSlot(c echo.Context) error {
    log.Println("Received request to update slot")

    ctx := c.Request().Context()
    slotId := c.Param("slot_id")

    req := new(booking.Slot)
    if err := c.Bind(req); err != nil {
        return response.ErrResponse(c, http.StatusBadRequest, "Invalid request payload")
    }

    // Since req.StartTime and req.EndTime are strings now (in "HH:mm" format), you can use them directly
    startTimeStr := req.StartTime
    endTimeStr := req.EndTime

    // Validate that startTime and endTime are not empty
    if startTimeStr == "" || endTimeStr == "" {
        return response.ErrResponse(c, http.StatusBadRequest, "Start time or end time cannot be empty")
    }

    // Call the usecase to update the slot with string times
    slot, err := h.slotsUsecase.UpdateSlot(ctx, slotId, startTimeStr, endTimeStr)
    if err != nil {
        return response.ErrResponse(c, http.StatusBadRequest, err.Error())
    }

    return response.SuccessResponse(c, http.StatusOK, slot)
}


func (h *slotsHttpHandler) FindSlot(c echo.Context) error {
	log.Println("Received request to find slot")

	ctx := c.Request().Context()
	slotId := c.Param("slot_id")

	slot, err := h.slotsUsecase.FindOneSlot(ctx, slotId)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, slot)
}


func (h *slotsHttpHandler) FindAllSlots(c echo.Context) error {
	log.Println("Received request to find all slots")

	ctx := c.Request().Context()

	slots, err := h.slotsUsecase.FindAllSlots(ctx)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, slots)
}

func (h *slotsHttpHandler) EnableOrDisableSlot(c echo.Context) error {
	log.Println("Received request to enable or disable slot")

	ctx := c.Request().Context()
	slotId := c.Param("slot_id")
	status := c.QueryParam("status")

	statusInt := 0
	if status == "1" {
		statusInt = 1
	}

	slot, err := h.slotsUsecase.EnableOrDisableSlot(ctx, slotId, statusInt)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, slot)
}