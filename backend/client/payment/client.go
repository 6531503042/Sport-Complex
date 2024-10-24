package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type PaymentClient struct {
	baseURL string
}

func NewPaymentClient(baseURL string) *PaymentClient {
	return &PaymentClient{baseURL: baseURL}
}

type CreatePaymentRequest struct {
	Amount        float64 `json:"amount" `                        // ต้องมี
	UserID        string  `json:"user_id" validate:"required"`    // ต้องมี
	BookingID     string  `json:"booking_id" validate:"required"` // ต้องมี
	PaymentMethod string  `json:"payment_method"`                 // สามารถไม่มีได้
	Currency      string  `json:"currency" validate:"required"`   // ต้องมี
}

type PaymentResponse struct {
	// ใส่ฟิลด์ที่ต้องการใน response
	ID        string  `json:"id"`
	Status    string  `json:"status"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	BookingID string  `json:"booking_id"`
}

func (c *PaymentClient) CreatePayment(req CreatePaymentRequest) (*PaymentResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(c.baseURL+"/payments", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to create payment: %s", resp.Status)
	}

	var paymentResp PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&paymentResp); err != nil {
		return nil, err
	}

	return &paymentResp, nil
}
