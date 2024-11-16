package handler

import (
	"context"
	"fmt"
	"main/modules/booking"
	bookingPb "main/modules/booking/proto"
	"main/modules/booking/usecase"
	"time"
)

type bookingGrpcHandler struct {
    bookingUsecase usecase.BookingUsecaseService
    bookingPb.UnimplementedBookingServiceServer
}

func NewBookingGrpcHandler(bookingUsecase usecase.BookingUsecaseService) bookingPb.BookingServiceServer {
    return &bookingGrpcHandler{
        bookingUsecase: bookingUsecase,
    }
}

func (h *bookingGrpcHandler) CreateBooking(ctx context.Context, req *bookingPb.CreateBookingRequest) (*bookingPb.BookingResponse, error) {
    bookingReq := &booking.CreateBookingRequest{
        UserId:          req.UserId,
        SlotType:        req.SlotType,
    }

    if req.SlotId != "" {
        bookingReq.SlotId = &req.SlotId
    }
    if req.BadmintonSlotId != "" {
        bookingReq.BadmintonSlotId = &req.BadmintonSlotId
    }

    result, err := h.bookingUsecase.InsertBooking(ctx, req.FacilityName, bookingReq)
    if err != nil {
        return &bookingPb.BookingResponse{
            ErrorMessage: fmt.Sprintf("Failed to create booking: %v", err),
        }, nil
    }

    return convertToBookingPbResponse(result), nil
}

func (h *bookingGrpcHandler) GetBooking(ctx context.Context, req *bookingPb.GetBookingRequest) (*bookingPb.BookingResponse, error) {
    result, err := h.bookingUsecase.FindBooking(ctx, req.BookingId)
    if err != nil {
        return &bookingPb.BookingResponse{
            ErrorMessage: fmt.Sprintf("Failed to find booking: %v", err),
        }, nil
    }

    bookingResponse := &booking.BookingResponse{
        Id:              result.Id,
        UserId:          result.UserId,
        SlotId:          result.SlotId,
        BadmintonSlotId: result.BadmintonSlotId,
        SlotType:        result.SlotType,
        Status:          result.Status,
        PaymentID:       result.PaymentID,
        QRCodeURL:       result.QRCodeURL,
        CreatedAt:       result.CreatedAt,
        UpdatedAt:       result.UpdatedAt,
    }

    return convertToBookingPbResponse(bookingResponse), nil
}


func (h *bookingGrpcHandler) UpdateBookingStatus(ctx context.Context, req *bookingPb.UpdateBookingStatusRequest) (*bookingPb.BookingResponse, error) {
    result, err := h.bookingUsecase.UpdateBooking(ctx, req.BookingId, req.Status)
    if err != nil {
        return &bookingPb.BookingResponse{
            ErrorMessage: fmt.Sprintf("Failed to update booking status: %v", err),
        }, nil
    }

    bookingResponse := &booking.BookingResponse{
        Id:              result.Id,
        UserId:          result.UserId,
        SlotId:          result.SlotId,
        BadmintonSlotId: result.BadmintonSlotId,
        SlotType:        result.SlotType,
        Status:          result.Status,
		PaymentID:       result.PaymentID,
        QRCodeURL:       result.QRCodeURL,
        CreatedAt:       result.CreatedAt,
        UpdatedAt:       result.UpdatedAt,
    }

    return convertToBookingPbResponse(bookingResponse), nil
}

func (h *bookingGrpcHandler) GetUserBookings(ctx context.Context, req *bookingPb.GetUserBookingsRequest) (*bookingPb.GetUserBookingsResponse, error) {
    bookings, err := h.bookingUsecase.FindOneUserBooking(ctx, req.UserId)
    if err != nil {
        return &bookingPb.GetUserBookingsResponse{
            ErrorMessage: fmt.Sprintf("Failed to find user bookings: %v", err),
        }, nil
    }

    var pbBookings []*bookingPb.BookingResponse
    for _, b := range bookings {
        // Convert booking.Booking to booking.BookingResponse
        bookingResponse := &booking.BookingResponse{
            Id:              b.Id,
            UserId:          b.UserId,
            SlotId:          b.SlotId,
            BadmintonSlotId: b.BadmintonSlotId,
            SlotType:        b.SlotType,
            Status:          b.Status,
            PaymentID:       b.PaymentID,
            QRCodeURL:       b.QRCodeURL,
            CreatedAt:       b.CreatedAt,
            UpdatedAt:       b.UpdatedAt,
        }
        pbBookings = append(pbBookings, convertToBookingPbResponse(bookingResponse))
    }
    return &bookingPb.GetUserBookingsResponse{
        Bookings: pbBookings,
    }, nil
}


func convertToBookingPbResponse(b *booking.BookingResponse) *bookingPb.BookingResponse {
    if b == nil {
        return &bookingPb.BookingResponse{}
    }

    var slotId, badmintonSlotId string
    if b.SlotId != nil {
        slotId = *b.SlotId
    }
    if b.BadmintonSlotId != nil {
        badmintonSlotId = *b.BadmintonSlotId
    }

    return &bookingPb.BookingResponse{
        Id:              b.Id.Hex(),
        UserId:          b.UserId,
        SlotId:          slotId,
        BadmintonSlotId: badmintonSlotId,
        SlotType:        b.SlotType,
        Status:          b.Status,
        PaymentId:       b.PaymentID,
        QrCodeUrl:       b.QRCodeURL,
        CreatedAt:       b.CreatedAt.Format(time.RFC3339),
        UpdatedAt:       b.UpdatedAt.Format(time.RFC3339),
    }
} 