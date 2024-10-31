package payment

import "time"

type CreatePaymentRequest struct {
	UserId        string  `json:"user_id" validate:"required"`        // ID of the user making the payment
	BookingId     string  `json:"booking_id" validate:"required"`     // ID of the booking related to the payment
	Amount        float64 `json:"amount" validate:"required,gt=0"`    // Amount to be paid
	Currency      string  `json:"currency" validate:"required"`        // Currency of the payment, e.g., THB, USD
	PaymentMethod string  `json:"payment_method" validate:"required"` // Method of payment, e.g., PromptPay
}

// PaymentRequestModel ใช้สำหรับรับข้อมูลการสร้าง payment จาก API

type PaymentRequest struct {
    UserId        string  `json:"user_id" validate:"required"`        // ID ของผู้ใช้ที่ทำการชำระเงิน
    BookingId     string  `json:"booking_id" validate:"required"`     // ID ของการจอง
    Amount        float64 `json:"amount" validate:"required"`         // จำนวนเงิน
    Currency      string  `json:"currency" validate:"required"`       // สกุลเงิน เช่น THB, USD
    PaymentMethod string  `json:"payment_method" validate:"required"` // วิธีการชำระเงิน เช่น PromptPay, CreditCard
}

// PaymentResponseModel ใช้สำหรับส่งข้อมูลการชำระเงินกลับไปยัง client
type PaymentResponse struct {
    Id            string        `json:"id"`             // ID ของการชำระเงิน
    PaymentID     string        `json:"payment_id"`
    UserId        string        `bson:"user_id" json:"user_id"`        // ID ของผู้ใช้ที่ทำการชำระเงิน
    BookingId     string        `json:"booking_id"`     // ID ของการจอง
    Amount        float64       `json:"amount"`         // จำนวนเงินที่ชำระ
    Currency      string        `json:"currency"`       // สกุลเงินที่ใช้
    PaymentMethod string        `json:"payment_method"` // วิธีการชำระเงิน
    QRCodeURL     string        `json:"qr_code_url"`    // URL ของ QR Code สำหรับการชำระเงิน
    Status        PaymentStatus `json:"status"`         // สถานะการชำระเงิน
    CreatedAt     time.Time     `json:"created_at"`     // เวลาที่สร้าง
    UpdatedAt     time.Time     `json:"updated_at"`     // เวลาที่อัปเดตล่าสุด
}

// NewPaymentResponseModel แปลง PaymentEntity ให้เป็น PaymentResponseModel
func NewPaymentResponse(payment *PaymentEntity) *PaymentResponse {
    return &PaymentResponse{
        Id:            payment.Id.Hex(),
        PaymentID:     payment.PaymentID,  // ส่งค่า payment_id กลับไปใน response ด้วย
        UserId:        payment.UserID,
        BookingId:     payment.BookingID,
        Amount:        payment.Amount,
        Currency:      payment.Currency,
        PaymentMethod: payment.PaymentMethod,
        QRCodeURL:     payment.QRCodeURL,
        Status:        payment.Status,
        CreatedAt:     payment.CreatedAt,
        UpdatedAt:     payment.UpdatedAt,
    }
}
