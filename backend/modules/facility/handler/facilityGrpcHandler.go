package handler

import (
	"context"
	"fmt"
	facilityPb "main/modules/facility/proto"
	"main/modules/facility/usecase"
)

type facilityGrpcHandler struct {
    facilityUsecase usecase.FacilityUsecaseService
    facilityPb.UnimplementedFacilityServiceServer
}

func NewFacilityGrpcHandler(facilityUsecase usecase.FacilityUsecaseService) *facilityGrpcHandler {
    return &facilityGrpcHandler{
        facilityUsecase: facilityUsecase,
    }
}

func (h *facilityGrpcHandler) CheckSlotAvailability(ctx context.Context, req *facilityPb.CheckSlotRequest) (*facilityPb.SlotAvailabilityResponse, error) {
    slot, err := h.facilityUsecase.FindOneSlot(ctx, req.FacilityName, req.SlotId)
    if err != nil {
        return &facilityPb.SlotAvailabilityResponse{
            IsAvailable:  false,
            ErrorMessage: fmt.Sprintf("Failed to find slot: %v", err),
        }, nil
    }

    isAvailable := slot.CurrentBookings < slot.MaxBookings
    return &facilityPb.SlotAvailabilityResponse{
        IsAvailable:     isAvailable,
        CurrentBookings: int32(slot.CurrentBookings),
        MaxBookings:     int32(slot.MaxBookings),
    }, nil
}

func (h *facilityGrpcHandler) GetFacilityPrice(ctx context.Context, req *facilityPb.FacilityPriceRequest) (*facilityPb.FacilityPriceResponse, error) {
    facility, err := h.facilityUsecase.FindOneFacility(ctx, "", req.FacilityName)
    if err != nil {
        return &facilityPb.FacilityPriceResponse{
            ErrorMessage: fmt.Sprintf("Failed to find facility: %v", err),
        }, nil
    }

    var price float64
    if req.UserType == "insider" {
        price = facility.PriceInsider
    } else {
        price = facility.PriceOutsider
    }

    return &facilityPb.FacilityPriceResponse{
        Price:    price,
        Currency: "THB",
    }, nil
}

func (h *facilityGrpcHandler) UpdateSlotBookingCount(ctx context.Context, req *facilityPb.UpdateSlotRequest) (*facilityPb.UpdateSlotResponse, error) {
    slot, err := h.facilityUsecase.FindOneSlot(ctx, req.FacilityName, req.SlotId)
    if err != nil {
        return &facilityPb.UpdateSlotResponse{
            Success:      false,
            ErrorMessage: fmt.Sprintf("Failed to find slot: %v", err),
        }, nil
    }

    // Update current bookings
    slot.CurrentBookings += int(req.Increment)
    
    // Ensure we don't go below 0 or above max
    if slot.CurrentBookings < 0 {
        slot.CurrentBookings = 0
    } else if slot.CurrentBookings > slot.MaxBookings {
        return &facilityPb.UpdateSlotResponse{
            Success:      false,
            ErrorMessage: "Booking count would exceed maximum allowed",
        }, nil
    }

    // Update the slot
    _, err = h.facilityUsecase.UpdateSlot(ctx, req.FacilityName, slot)
    if err != nil {
        return &facilityPb.UpdateSlotResponse{
            Success:      false,
            ErrorMessage: fmt.Sprintf("Failed to update slot: %v", err),
        }, nil
    }

    return &facilityPb.UpdateSlotResponse{
        Success: true,
    }, nil
} 