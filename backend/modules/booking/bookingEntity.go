package booking

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Booking struct {
		Id              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
		UserId          string             `bson:"user_id" json:"user_id"`
		SlotId          *string            `bson:"slot_id,omitempty" json:"slot_id,omitempty"`                     // String type for normal slot ID
		BadmintonSlotId *string            `bson:"badminton_slot_id,omitempty" json:"badminton_slot_id,omitempty"` // String type for badminton slot ID
		SlotType        string             `bson:"slot_type" json:"slot_type"`       // "normal" or "badminton"
		Status          string             `bson:"status" json:"status"`
		PaymentId       string            `bson:"payment_id,omitempty" json:"payment_id,omitempty"`
		Facility        string            `bson:"facility" json:"facility"`
		CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
		UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
	}
)
