package facility

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Slot struct {
		Id              primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
		StartTime       string             `bson:"start_time" json:"start_time"`
		EndTime         string             `bson:"end_time" json:"end_time"`
		Status          int                `bson:"status" json:"status"`
		MaxBookings     int                `bson:"max_bookings" json:"max_bookings"`
		CurrentBookings int                `bson:"current_bookings" json:"current_bookings"`
		FacilityType    string             `bson:"facility_type" json:"facility_type"` // e.g., "gym", "swimming", "football"
		CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
		UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
	}

	BadmintonCourt struct {
		Id          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
		CourtNumber int                `bson:"court_number" json:"court_number"` // Court number, e.g., 1, 2, 3, 4
		Status      int                `bson:"status" json:"status"`             // e.g., 0 = available, 1 = booked
	}

	BadmintonSlot struct {
		Id              primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
		StartTime       string            `bson:"start_time" json:"start_time"`
		EndTime         string            `bson:"end_time" json:"end_time"`
		CourtId         primitive.ObjectID `bson:"court_id" json:"court_id"`
		CourtNumber    int                `bson:"court_number"`
		Status          int               `bson:"status" json:"status"`
		MaxBookings     int               `bson:"max_bookings" json:"max_bookings"`
		CurrentBookings int               `bson:"current_bookings" json:"current_bookings"`
		CreatedAt       time.Time         `bson:"created_at" json:"created_at"`
		UpdatedAt       time.Time         `bson:"updated_at" json:"updated_at"`
	}
)