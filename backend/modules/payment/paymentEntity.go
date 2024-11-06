package payment

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PaymentStatus ชนิดของสถานะการชำระเงิน
type PaymentStatus string

const (
	Pending   PaymentStatus = "PENDING"
	Completed PaymentStatus = "COMPLETED"
	Failed    PaymentStatus = "FAILED"
	Canceled  PaymentStatus = "CANCELED"
)

// PaymentEntity เป็นโครงสร้างข้อมูลสำหรับการจัดเก็บ transaction การชำระเงิน
type PaymentEntity struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id"`              // MongoDB ObjectID
	UserID        string             `bson:"user_id" json:"user_id"`                // ID of the user making the payment
	BookingID     string             `bson:"booking_id" json:"booking_id"`          // ID of the booking associated with the payment
	PaymentID     string             `bson:"payment_id" json:"payment_id"`
	Amount        float64            `bson:"amount" json:"amount"`                  // Amount to be paid
	Currency      string             `bson:"currency" json:"currency"`              // Currency used, e.g., THB, USD
	PaymentMethod string             `bson:"payment_method" json:"payment_method"`  // Payment method, e.g., PromptPay, CreditCard
	QRCodeURL     string             `bson:"qr_code_url" json:"qr_code_url"`       // URL of the QR Code for payment
	Status        PaymentStatus      `bson:"status" json:"status"`                  // Payment status (Pending, Completed, Failed)
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`          // Time when the payment record was created
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`          // Time when the record was last updated
}

type PaymentSlip struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        string    `bson:"user_id"`
	BookingID     string    `bson:"booking_id"`
    ImagePath     string    `bson:"image_path"`
    Status        string    `bson:"status"`  // e.g., Pending, Approved, Rejected
    SubmittedDate time.Time `bson:"submitted_date"`
}
