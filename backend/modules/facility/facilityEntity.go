package facility

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Facilitiy struct {
		Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
		Name        string             `bson:"name" json:"name"`
		PriceInsider float64           `bson:"price_insider" json:"price_insider"`
		PriceOutsider float64          `bson:"price_outsider" json:"price_outsider"`
		Description string             `bson:"description" json:"description"`
		CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
		UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	}

	FacilityBson struct {
		Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
		Name        string             `bson:"name" json:"name"`
		PriceInsider float64           `bson:"price_insider" json:"price_insider"`
		PriceOutsider float64          `bson:"price_outsider" json:"price_outsider"`
		Description string             `bson:"description" json:"description"`
		CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
		UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	}

	Slot struct {
		Id        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
		StartTime string             `bson:"start_time" json:"start_time"`
		EndTime   string             `bson:"end_time" json:"end_time"`
		Status    int                `bson:"status" json:"status"`
		CreatedAt time.Time          `bson:"created_at" json:"created_at"`
		UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	}
)