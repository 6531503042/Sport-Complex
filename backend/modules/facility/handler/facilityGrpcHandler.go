package handler

import (
	"context"
	"fmt"
	"main/modules/facility"
	facilityPb "main/modules/facility/proto"
	"main/modules/facility/usecase"
)

type facilityGrpcHandler struct {
    facilityUsecase usecase.FacilityUsecaseService
    facilityPb.UnimplementedFacilityServiceServer
}

func NewFacilityGrpcHandler(facilityUsecase usecase.FacilityUsecaseService) facilityPb.FacilityServiceServer {
    return &facilityGrpcHandler{
        facilityUsecase: facilityUsecase,
    }
}

func (h *facilityGrpcHandler) CheckSlotAvailability(ctx context.Context, req *facilityPb.CheckSlotRequest) (*facilityPb.SlotAvailabilityResponse, error) {
    var slot interface{}
    var err error

    if req.SlotType == "badminton" {
        // Handle badminton slot check
        slots, err := h.facilityUsecase.FindBadmintonSlot(ctx)
        if err != nil {
            return &facilityPb.SlotAvailabilityResponse{
                IsAvailable:  false,
                ErrorMessage: fmt.Sprintf("Failed to find badminton slot: %v", err),
            }, nil
        }
        // Find the specific slot
        for _, s := range slots {
            if s.Id.Hex() == req.SlotId {
                return &facilityPb.SlotAvailabilityResponse{
                    IsAvailable: s.Status == 0,
                }, nil
            }
        }
    } else {
        // Handle normal slot check
        slot, err = h.facilityUsecase.FindOneSlot(ctx, req.FacilityName, req.SlotId)
        if err != nil {
            return &facilityPb.SlotAvailabilityResponse{
                IsAvailable:  false,
                ErrorMessage: fmt.Sprintf("Failed to find slot: %v", err),
            }, nil
        }
        if s, ok := slot.(*facility.Slot); ok {
            return &facilityPb.SlotAvailabilityResponse{
                IsAvailable:     s.CurrentBookings < s.MaxBookings,
                CurrentBookings: int32(s.CurrentBookings),
                MaxBookings:     int32(s.MaxBookings),
            }, nil
        }
    }

    return &facilityPb.SlotAvailabilityResponse{
        IsAvailable:  false,
        ErrorMessage: "Slot not found",
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

    slot.CurrentBookings += int(req.Increment)
    
    if slot.CurrentBookings < 0 {
        slot.CurrentBookings = 0
    } else if slot.CurrentBookings > slot.MaxBookings {
        return &facilityPb.UpdateSlotResponse{
            Success:      false,
            ErrorMessage: "Booking count would exceed maximum allowed",
        }, nil
    }

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