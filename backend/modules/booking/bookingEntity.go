package booking

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Booking struct {
		Id        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
		UserID    string             `json:"user_id" bson:"user_id"`
		GymId	string 				`json:"gym_id", bson:"gym_id"`
		SwimmingId	string			`json:"swimming_id", bson:"swimming_id"`
		BadmintonId 	string		`json:"badminton_id", bson:"badminton_id"`
		FootballId		string		`json:"football_id", bson:"foorball_id"`
		Slot	TimeSlot		`json:"slot", bson:"slot"`
		Status 	string		`json:"status", bson:"status"`
		CreatedAt time.Time          `json:"created_at" bson:"created_at"`
		UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	}

	TimeSlot struct {
		Id        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
		Start 	string		`json:"start", bson:"start"`
		End 		string		`json:"end", bson:"end"`
		Status 	string		`json:"status", bson:"status"`
	}
)