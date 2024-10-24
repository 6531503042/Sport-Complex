package booking

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type(
	BookingDetails struct {
		Id        string    `json:"id"`
		UserId    string    `json:"user_id"`
		GymId     string    `json:"gym_id"`
		SwimmingId     string    `json:"swimming_id"`
		BadmintonId    string    `json:"badminton_id"`
		FootballId     string    `json:"football_id"`
		Slot      TimeSlot  `json:"slot"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	CreateBookingReq struct {
		UserId    string    `json:"user_id"`
		GymId     string    `json:"gym_id"`
		SwimmingId     string    `json:"swimming_id"`
		BadmintonId    string    `json:"badminton_id"`
		FootballId     string    `json:"football_id"`
		Slot      TimeSlot  `json:"slot"`
	}
)
