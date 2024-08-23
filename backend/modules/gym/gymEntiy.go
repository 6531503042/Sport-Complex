package gym

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (

	GymSlot struct {
		Id	primitive.ObjectID	`json:"_id" bson:"_id,omitempty"`
		StartTime string `json:"start_time", bson:"start_time"`
		EndTime string `json:"end_time", bson:"end_time"`
		MaxQueue int `json:"max_queue", bson:"max_queue"`
		Booked int `json:"booked", bson:"booked"`
		IsClose bool `json:"is_closed", bson:"is_closed"`
	}

	
)