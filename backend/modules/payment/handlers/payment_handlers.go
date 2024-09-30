package payment

import (
	"context"
	"fmt"
	"log"
)

type PaymentServiceServer struct {
	UnimplementedPaymentServiceServer
}

// ฟังก์ชันสำหรับการสร้างการชำระเงินและคืน QR Code
func (s *PaymentServiceServer) CreatePayment(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error) {
	log.Printf("Received CreatePayment request for user_id: %s, amount: %.2f", req.UserId, req.Amount)
	qrCode := generateQRCode(req.Amount)
	return &PaymentResponse{QrCode: qrCode}, nil
}

// ฟังก์ชันสำหรับตรวจสอบสถานะการชำระเงิน
func (s *PaymentServiceServer) CheckPaymentStatus(ctx context.Context, req *PaymentStatusRequest) (*PaymentStatusResponse, error) {
	log.Printf("Checking payment status for user_id: %s", req.UserId)
	status := checkPaymentStatus(req.UserId)
	return &PaymentStatusResponse{Status: status}, nil
}

// ฟังก์ชันสร้าง QR Code (จำลอง)
func generateQRCode(amount float32) string {
	return fmt.Sprintf("QR_CODE_FOR_AMOUNT_%.2f", amount)
}

// ฟังก์ชันตรวจสอบสถานะการชำระเงิน (จำลอง)
func checkPaymentStatus(userId string) string {
	// สามารถเพิ่มการตรวจสอบจากฐานข้อมูลได้
	return "PAID" // จำลองว่าการชำระเงินเสร็จสิ้นแล้ว
}
