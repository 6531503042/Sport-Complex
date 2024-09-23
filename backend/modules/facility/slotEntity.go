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
		Id              primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
		CourtNumber     int                `bson:"court_number" json:"court_number"` // Court number, e.g., 1, 2, 3
		IsBooked        bool               `bson:"is_booked" json:"is_booked"`       // True if the court is booked
	}

	BadmintonSlot struct {
		Id              primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
		StartTime       string             `bson:"start_time" json:"start_time"`
		EndTime         string             `bson:"end_time" json:"end_time"`
		CourtId 		string				`bson:"court_id" json:"court_id"`
		// Courts          []BadmintonCourt   `bson:"courts" json:"courts"`
		Status          int                `bson:"status" json:"status"`
		CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
		UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
	}
)