package usecase

import (
	"context"
	"fmt"
	"log"
	"main/config"
	"main/modules/payment"
	"main/modules/payment/repository"

	"net/url"
	"time"

	"github.com/Frontware/promptpay"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentUsecaseService interface {
	CreatePayment(ctx context.Context, userId, bookingId, paymentMethod, facilityName string, amount float64) (*payment.PaymentResponse, error)
	UpdatePayment(ctx context.Context, paymentId, status string) (*payment.PaymentEntity, error)
	FindPayment(ctx context.Context, paymentId string) (*payment.PaymentEntity, error)
	FindPaymentsByUser(ctx context.Context, userId string) ([]payment.PaymentEntity, error)
	SaveSlip(ctx context.Context, slip payment.PaymentSlip) error
	FindSlipByUserId(ctx context.Context, userId string) ([]payment.PaymentSlip, error)
	UpdateSlipStatus(ctx context.Context, slipId string, newStatus string) error
	GetPendingSlips(ctx context.Context) ([]payment.PaymentSlip, error)
}

type paymentUsecase struct {
	cfg               *config.Config
	paymentRepository repository.PaymentRepositoryService
}

// NewPaymentUsecase creates and returns a new payment usecase instance.
func NewPaymentUsecase(cfg *config.Config, paymentRepository repository.PaymentRepositoryService) PaymentUsecaseService {
	return &paymentUsecase{
		cfg:               cfg,
		paymentRepository: paymentRepository,
	}
}

func (u *paymentUsecase) CreatePayment(ctx context.Context, userId, bookingId, paymentMethod, facilityName string, amount float64) (*payment.PaymentResponse, error) {
	paymentDoc := &payment.PaymentEntity{
		Id:            primitive.NewObjectID(),
		UserID:        userId,
		BookingID:     bookingId,
		Amount:        amount,
		Currency:      "THB",
		PaymentMethod: paymentMethod,
		FacilityName:  facilityName,
		Status:        payment.Pending,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Generate QR code (PromptPay logic)
	promptPay := &promptpay.PromptPay{
		PromptPayID: "1579901028845", // Replace with actual PromptPay ID
		Amount:      amount,
		OneTime:     true,
	}

	qrCodeData, err := promptPay.Gen()
	if err != nil {
		return nil, fmt.Errorf("error generating QR code: %w", err)
	}

	// Create the QR code URL
	qrCodeURL := fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?data=%s&size=300x300", url.QueryEscape(qrCodeData))
	paymentDoc.QRCodeURL = qrCodeURL

	// Save payment in the repository
	paymentResult, err := u.paymentRepository.InsertPayment(ctx, paymentDoc)
	if err != nil {
		return nil, fmt.Errorf("error creating payment: %w", err)
	}

	// Convert the entity to response
	response := &payment.PaymentResponse{
		Id:           paymentResult.Id.Hex(),
		PaymentID:    paymentResult.Id.Hex(),
		UserId:       paymentResult.UserID,
		BookingId:    paymentResult.BookingID,
		Amount:       paymentResult.Amount,
		Currency:     paymentResult.Currency,
		PaymentMethod: paymentResult.PaymentMethod,
		Status:       paymentResult.Status,
		FacilityName: paymentResult.FacilityName,
		CreatedAt:    paymentResult.CreatedAt,
		UpdatedAt:    paymentResult.UpdatedAt,
		QRCodeURL:    paymentResult.QRCodeURL,
	}

	return response, nil
}

func (u *paymentUsecase) UpdatePayment(ctx context.Context, paymentId, status string) (*payment.PaymentEntity, error) {
	paymentEntity, err := u.paymentRepository.FindPayment(ctx, paymentId)
	if err != nil {
		return nil, fmt.Errorf("failed to find payment: %w", err)
	}

	paymentEntity.Status = payment.PaymentStatus(status)
	paymentEntity.UpdatedAt = time.Now()

	// Update the payment status
	updatedPayment, err := u.paymentRepository.UpdatePayment(ctx, paymentEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	return updatedPayment, nil
}

func (u *paymentUsecase) FindPayment(ctx context.Context, paymentId string) (*payment.PaymentEntity, error) {
	return u.paymentRepository.FindPayment(ctx, paymentId)
}

func (u *paymentUsecase) FindPaymentsByUser(ctx context.Context, userId string) ([]payment.PaymentEntity, error) {
	// Call the repository to retrieve payments
	payments, err := u.paymentRepository.FindPaymentsByUser(ctx, userId)
	if err != nil {
		log.Printf("Error: FindPaymentsByUser failed: %s", err.Error())
		return nil, fmt.Errorf("failed to find payments by user ID: %w", err)
	}

	return payments, nil
}


func (u *paymentUsecase) SaveSlip(ctx context.Context, slip payment.PaymentSlip) error {
	err := u.paymentRepository.SaveSlip(ctx, slip)
	if err != nil {
		return fmt.Errorf("error saving slip: %w", err)
	}
	return nil
}

func (u *paymentUsecase) FindSlipByUserId(ctx context.Context, userId string) ([]payment.PaymentSlip, error) {
	slips, err := u.paymentRepository.FindSlipByUserId(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("error finding slips by user ID: %w", err)
	}
	return slips, nil
}

func (u *paymentUsecase) UpdateSlipStatus(ctx context.Context, slipId string, newStatus string) error {
	err := u.paymentRepository.UpdateSlipStatus(ctx, slipId, newStatus)
	if err != nil {
		return fmt.Errorf("error updating slip status: %w", err)
	}
	return nil
}

func (u *paymentUsecase) GetPendingSlips(ctx context.Context) ([]payment.PaymentSlip, error) {
	slips, err := u.paymentRepository.GetPendingSlips(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving pending slips: %w", err)
	}
	return slips, nil
}