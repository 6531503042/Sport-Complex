package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PaymentClient struct {
	baseURL string
}

func NewPaymentClient(baseURL string) *PaymentClient {
	return &PaymentClient{baseURL: baseURL}
}

type CreatePaymentRequest struct {
    Amount        float64 `json:"amount"`
    UserID        string  `json:"user_id"`
    BookingID     string  `json:"booking_id"`
    PaymentMethod string  `json:"payment_method"`
    Currency      string  `json:"currency"`
    FacilityName  string  `json:"facility_name"` // เพิ่มฟิลด์นี้
}


type PaymentResponse struct {
	ID        string  `json:"id"`
	Status    string  `json:"status"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	BookingID string  `json:"booking_id"`
	Message   string  `json:"message,omitempty"`
	QRCodeURL string  `json:"qr_code_url,omitempty"` // Add this field to return the QR code URL
}


func (c *PaymentClient) CreatePayment(req CreatePaymentRequest) (*PaymentResponse, error) {
	// Handle specific logic for PromptPay
	if req.PaymentMethod == "PROMPTPAY" {
		// Add logic to generate PromptPay QR code
		body, err := json.Marshal(req)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}

		// Send the request to your service to create a PromptPay payment
		resp, err := http.Post(c.baseURL+"/payments/promptpay", "application/json", bytes.NewBuffer(body))
		if err != nil {
			return nil, fmt.Errorf("failed to send request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := ioutil.ReadAll(resp.Body)
			return nil, fmt.Errorf("failed to create payment, status: %s, response: %s", resp.Status, string(respBody))
		}

		var paymentResp PaymentResponse
		if err := json.NewDecoder(resp.Body).Decode(&paymentResp); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		// Return the response, which should include the QR code URL
		return &paymentResp, nil
	}

	// For other payment methods
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(c.baseURL+"/payments", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create payment, status: %s, response: %s", resp.Status, string(respBody))
	}

	var paymentResp PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&paymentResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &paymentResp, nil
}


// CheckPaymentStatusRequest ฟังก์ชันสำหรับรับ request การตรวจสอบสถานะการชำระเงิน
type CheckPaymentStatusRequest struct {
	PaymentID string `json:"payment_id" validate:"required"`
}

// CheckPaymentStatus ฟังก์ชันสำหรับตรวจสอบสถานะการชำระเงิน
func (c *PaymentClient) CheckPaymentStatus(paymentID string) (*PaymentResponse, error) {
	url := fmt.Sprintf("%s/payments/promptpay/%s/status", c.baseURL, paymentID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to check payment status, status: %s, response: %s", resp.Status, string(respBody))
	}

	var paymentResp PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&paymentResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &paymentResp, nil
}


func (c *PaymentClient) UpdatePaymentStatus(paymentID string, status string) error {
	updateRequest := struct {
		Status string `json:"status"`
	}{
		Status: status,
	}

	body, err := json.Marshal(updateRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/payments/%s", c.baseURL, paymentID), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update payment status: %s", resp.Status)
	}

	return nil
}

func (c *PaymentClient) UpdatePaymentToCompleted(paymentID string) error {
    return c.UpdatePaymentStatus(paymentID, "COMPLETED")
}


