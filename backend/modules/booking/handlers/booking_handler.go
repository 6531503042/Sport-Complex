package handlers

import (
	"context"
	"fmt"
	"log"

	"main/modules/booking/kafka"
	"main/modules/booking/models"
	"main/modules/booking/repository"
	pb "main/modules/booking/proto" // Correct import for protobuf
)

type BookingHandler struct {
	repo   repository.BookingRepository
	kafka  *kafka.KafkaProducer
	pb.UnimplementedBookingServiceServer
}

func NewBookingHandler(repo repository.BookingRepository, kafka *kafka.KafkaProducer) *BookingHandler {
	return &BookingHandler{repo: repo, kafka: kafka}
}

func (h *BookingHandler) CreateBooking(ctx context.Context, req *pb.CreateBookingRequest) (*pb.CreateBookingResponse, error) {
	booking := models.Booking{
		UserID:     req.GetUserId(),
		FacilityID: req.GetFacilityId(),
		Timeslot:   req.GetTimeslot(),
		Price:      req.GetPrice(),
		Status:     "Pending",
	}

	err := h.repo.CreateBooking(ctx, &booking)
	if err != nil {
		return nil, fmt.Errorf("failed to create booking: %v", err)
	}

	// Trigger Kafka event for payment processing
	err = h.kafka.PublishBookingEvent(ctx, booking)
	if err != nil {
		log.Printf("failed to publish booking event: %v", err)
	}

	return &pb.CreateBookingResponse{
		BookingId: booking.ID.Hex(),
		Status:    booking.Status,
	}, nil
}

// Implement other methods for GetBookingDetails, ListBookings, CancelBooking, UpdateBookingStatus...
