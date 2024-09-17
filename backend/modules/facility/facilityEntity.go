package facility

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Facilitiy struct {
		ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
		Name        string             `bson:"name" json:"name"`
		PriceInsider float64           `bson:"price_insider" json:"price_insider"`
		PriceOutsider float64          `bson:"price_outsider" json:"price_outsider"`
		Description string             `bson:"description" json:"description"`
		CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
		UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	}

	FacilityBson struct {
		ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
		Name        string             `bson:"name" json:"name"`
		PriceInsider float64           `bson:"price_insider" json:"price_insider"`
		PriceOutsider float64          `bson:"price_outsider" json:"price_outsider"`
		Description string             `bson:"description" json:"description"`
		CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
		UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	}
)