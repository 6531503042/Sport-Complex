package booking

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	// CreateBookingRequest supports booking either a normal slot or a badminton slot
	CreateBookingRequest struct {
		UserId          string  `json:"user_id" validate:"required,max=64"`
		SlotId          *string `json:"slot_id,omitempty"`             // Normal slot ID (optional if badminton)
		BadmintonSlotId *string `json:"badminton_slot_id,omitempty"`   // Badminton slot ID (optional if normal)
		SlotType        string  `json:"slot_type" validate:"required"` // "normal" or "badminton"
	}

	// BookingSearchRequest for searching bookings by user, slot, or status
	BookingSearchRequest struct {
		UserId          string  `json:"user_id,omitempty"`
		SlotId          *string `json:"slot_id,omitempty"`             // For normal slots
		BadmintonSlotId *string `json:"badminton_slot_id,omitempty"`   // For badminton slots
		Status          int     `json:"status,omitempty"`
		Limit           int     `json:"limit,omitempty"`
		Offset          int     `json:"offset,omitempty"`
	}

	// BookingUpdateRequest for updating the status or rescheduling bookings
	BookingUpdateRequest struct {
		Status          int     `json:"status,omitempty"`             // The new status of the booking (e.g., confirmed, canceled)
		SlotId          *string `json:"slot_id,omitempty"`            // Slot ID for normal slot rescheduling
		BadmintonSlotId *string `json:"badminton_slot_id,omitempty"`  // Badminton-specific slot ID for rescheduling
		StartAt         string  `json:"start_at,omitempty"`           // New start time for rescheduling
		EndAt           string  `json:"end_at,omitempty"`             // New end time for rescheduling
	}

	// BookingResponse returns booking details, supporting both slot types
	BookingResponse struct {
		Id              primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
		UserId          string             `bson:"user_id" json:"user_id"`
		SlotId          *string            `bson:"slot_id,omitempty"`             // Slot ID for normal facilities
		BadmintonSlotId *string            `bson:"badminton_slot_id,omitempty"`   // Slot ID for badminton-specific bookings
		SlotType        string             `json:"slot_type"`                     // "normal" or "badminton"
		Status          int                `json:"status"`
		CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
		UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
	}

	// EnableOrDisableBookingRequest is used to enable or disable a booking
	EnableOrDisableBookingRequest struct {
		Status int `json:"status" validate:"required"`
	}
)
