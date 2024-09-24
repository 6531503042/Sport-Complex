package booking

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Booking struct {
		Id              primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
		UserId          primitive.ObjectID `bson:"user_id" json:"user_id"`
		SlotId          *primitive.ObjectID `bson:"slot_id,omitempty" json:"slot_id,omitempty"`           
		BadmintonSlotId *primitive.ObjectID `bson:"badminton_slot_id,omitempty" json:"badminton_slot_id,omitempty"`
		SlotType        string             `bson:"slot_type" json:"slot_type"`   
		Status          int                `bson:"status" json:"status"`
		CreatedAt       time.Time          `bson:"created_at" json:"created_at"`        
		UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
	}
)
