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
	Message   string  `json:"message,omitempty"`
	
}

func (c *PaymentClient) CreatePayment(req CreatePaymentRequest) (*PaymentResponse, error) {
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
	// สร้าง URL สำหรับตรวจสอบสถานะการชำระเงิน โดยใช้ paymentID
	url := fmt.Sprintf("%s/payments/%s/status", c.baseURL, paymentID)

	// ส่ง request ไปยัง payment service
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// ตรวจสอบสถานะการตอบกลับ
	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to check payment status, status: %s, response: %s", resp.Status, string(respBody))
	}

	// แปลง response body เป็น PaymentResponse
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


