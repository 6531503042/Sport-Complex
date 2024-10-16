package booking

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	CreateBookingReq struct {
		Userid string `json:"user_id" validate:"required,max=64"`
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
        // Add pagination fields if needed
        Limit  int `json:"limit,omitempty"`
        Offset int `json:"offset,omitempty"`
    }

	BookingUpdateReq struct {
        Status int `json:"status" validate:"required"`
    }

    EnableOrDisableBookingReq struct {
        Status int `json:"status" validate:"required"`
    }
)