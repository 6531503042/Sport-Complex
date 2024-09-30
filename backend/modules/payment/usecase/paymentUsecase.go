package usecase

import (
	"context"
	"fmt"
	"main/config"
	"main/modules/payment"
	"main/modules/payment/repository"
	"main/pkg/utils"

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

	// Create the payment object to pass to the repository
	paymentDoc := &payment.PaymentEntity{
		UserID:        userId,
		BookingID:     bookingId,
		Amount:        amount,
		Currency:      "THB", // Set your currency here
		PaymentMethod: "PromptPay", // Or any other payment method
		Status:        payment.Pending,
		CreatedAt:     utils.LocalTime(),
		UpdatedAt:     utils.LocalTime(),
	}

	// Insert the payment via repository
	paymentResult, err := u.paymentRepository.InsertPayment(ctx, paymentDoc)
	if err != nil {
		return nil, fmt.Errorf("error creating payment: %w", err)
	}

	// Prepare and return response
	response := &payment.PaymentResponse{
		Id:          paymentResult.Id.Hex(),
		UserId:      paymentResult.UserID,
		BookingId:   paymentResult.BookingID,
		Amount:      paymentResult.Amount,
		Currency:    paymentResult.Currency,
		Status:      paymentResult.Status,
		CreatedAt:   paymentResult.CreatedAt,
		UpdatedAt:   paymentResult.UpdatedAt,
	}

	return response, nil
}

func (u *paymentUsecase) UpdatePayment(ctx context.Context, paymentId string, status string) (*payment.PaymentEntity, error) {
	paymentEntity, err := u.paymentRepository.FindPayment(ctx, paymentId)
	if err != nil {
		return nil, fmt.Errorf("error: failed to find payment: %w", err)
	}

	// Update the payment status
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
