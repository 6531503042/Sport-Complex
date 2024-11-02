package usecase

import (
	"context"
	"fmt"
	"main/config"
	"main/modules/payment"
	"main/modules/payment/repository"
	"main/pkg/utils"
	"net/url"

	"github.com/Frontware/promptpay"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	PaymentUsecaseService interface {
		CreatePayment(ctx context.Context, userId string, bookingId string, amount float64) (*payment.PaymentResponse, error)
		UpdatePayment(ctx context.Context, paymentId string, status string) (*payment.PaymentEntity, error)
		FindPayment(ctx context.Context, paymentId string) (*payment.PaymentEntity, error)
		FindPaymentsByUser(ctx context.Context, userId string) ([]payment.PaymentEntity, error)
	}

	paymentUsecase struct {
		cfg                 *config.Config
		paymentRepository   repository.PaymentRepositoryService
	}
)

// NewPaymentUsecase returns a new instance of PaymentUsecaseService using the given payment repository.
func NewPaymentUsecase(paymentRepository repository.PaymentRepositoryService) PaymentUsecaseService {
	return &paymentUsecase{
		cfg:               &config.Config{},
		paymentRepository: paymentRepository,
	}
}

func (u *paymentUsecase) CreatePayment(ctx context.Context, userId string, bookingId string, amount float64) (*payment.PaymentResponse, error) {
	paymentDoc := &payment.PaymentEntity{
		PaymentID:     primitive.NewObjectID().Hex(),
		UserID:        userId,
		BookingID:     bookingId,
		Amount:        amount,
		Currency:      "THB",
		PaymentMethod: "PromptPay",
		Status:        payment.Pending,
		CreatedAt:     utils.LocalTime(),
		UpdatedAt:     utils.LocalTime(),
	}

	// Define the PromptPay phone number
	promptPayPhoneNumber := "0947044119" // Replace with the actual phone number

	// Create a new PromptPay instance with the phone number
	promptPay := &promptpay.PromptPay{
		PromptPayID: promptPayPhoneNumber,
		Amount:      amount,
		OneTime:     true,
	}

	// Generate the QR code data
	qrCodeData, err := promptPay.Gen()
	if err != nil {
		return nil, fmt.Errorf("error generating QR code: %w", err)
	}

	// Create URL for the QR code
	qrCodeURL := fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?data=%s&size=300x300", url.QueryEscape(qrCodeData))
	paymentDoc.QRCodeURL = qrCodeURL

	// Insert the payment via repository
	paymentResult, err := u.paymentRepository.InsertPayment(ctx, paymentDoc)
	if err != nil {
		return nil, fmt.Errorf("error creating payment: %w", err)
	}

	response := &payment.PaymentResponse{
		Id:        paymentResult.Id.Hex(),
		PaymentID: paymentResult.PaymentID,
		UserId:    paymentResult.UserID,
		BookingId: paymentResult.BookingID,
		Amount:    paymentResult.Amount,
		Currency:  paymentResult.Currency,
		Status:    paymentResult.Status,
		CreatedAt: paymentResult.CreatedAt,
		UpdatedAt: paymentResult.UpdatedAt,
		QRCodeURL: paymentResult.QRCodeURL,
	}

	return response, nil
}

func (u *paymentUsecase) UpdatePayment(ctx context.Context, paymentId string, status string) (*payment.PaymentEntity, error) {
	paymentEntity, err := u.paymentRepository.FindPayment(ctx, paymentId)
	if err != nil {
		return nil, fmt.Errorf("error: failed to find payment: %w", err)
	}

	paymentEntity.Status = payment.PaymentStatus(status)
	paymentEntity.UpdatedAt = utils.LocalTime()

	updatedPayment, err := u.paymentRepository.UpdatePayment(ctx, paymentEntity)
	if err != nil {
		return nil, fmt.Errorf("error: failed to update payment: %w", err)
	}

	return updatedPayment, nil
}

func (u *paymentUsecase) FindPayment(ctx context.Context, paymentId string) (*payment.PaymentEntity, error) {
	return u.paymentRepository.FindPayment(ctx, paymentId)
}

func (u *paymentUsecase) FindPaymentsByUser(ctx context.Context, userId string) ([]payment.PaymentEntity, error) {
	return u.paymentRepository.FindPaymentsByUser(ctx, userId)
}
