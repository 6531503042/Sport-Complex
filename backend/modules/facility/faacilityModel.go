package facility

import "time"

type (
	CreateFaciliityRequest struct {
		Name          string  `json:"name" validate:"required"`
		PriceInsider  float64 `json:"price_insider" validate:"required"`
		PriceOutsider float64 `json:"price_outsider" validate:"required"`
		Description   string  `json:"description" validate:"required"`
	}

	UpdateFaciltyRequest struct {
		Name          string  `json:"name,omitempty"`
		PriceInsider  float64 `json:"price_insider,omitempty"`
		PriceOutsider float64 `json:"price_outsider,omitempty"`
		Description   string  `json:"description,omitempty"`
	}

	FacilitiyResponse struct {
		Id            string    `json:"id"`
		Name          string    `json:"name"`
		PriceInsider  float64   `json:"price_insider"`
		PriceOutsider float64   `json:"price_outsider"`
		Description   string    `json:"description"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
	}
)