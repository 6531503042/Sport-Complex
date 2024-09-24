package booking

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	CreateBookingReq struct {
		UserId string `json:"user_id" validate:"required,max=64"`
		SlotId   string `json:"slot_id" validate:"required,max=64"`
	}

	BookingShowCase struct {
		Id        primitive.ObjectID `json:"id"`
        UserId    string             `json:"user_id"`
        SlotId    string             `json:"slot_id"`
        Status    int                `json:"status"`
        CreatedAt time.Time          `json:"created_at"`
        UpdatedAt time.Time          `json:"updated_at"`
	}

	BookingSearchReq struct {
        UserId string `json:"user_id,omitempty"`
        SlotId string `json:"slot_id,omitempty"`
        Status int    `json:"status,omitempty"`
        Limit  int `json:"limit,omitempty"`
        Offset int `json:"offset,omitempty"`
    }

	BookingUpdateReq struct {
        Status  string `json:"status,omitempty"`  // The new status of the booking (e.g., confirmed, canceled)
		SlotId  int64  `json:"slot_id,omitempty"` // The slot ID for rescheduling (if allowed)
		StartAt string `json:"start_at,omitempty"` // New start time for rescheduling
		EndAt   string `json:"end_at,omitempty"`   // New end time for rescheduling
	}

    EnableOrDisableBookingReq struct {
        Status int `json:"status" validate:"required"`
    }
)