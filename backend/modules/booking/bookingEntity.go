package booking

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Booking struct {
		Id        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
		UserId    string             `bson:"user_id" json:"user_id"`
		SlotId    string             `bson:"slot_id" json:"slot_id"`
		Status    int                `bson:"status" json:"status"`
		CreatedAt time.Time          `bson:"created_at" json:"created_at"`
		UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	}
)
