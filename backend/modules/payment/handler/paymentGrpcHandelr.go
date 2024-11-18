// paymentGrpcHandler.go
package handler

import (
	"context"
	"fmt"
	
	paymentPb "main/modules/payment/proto"
	"main/modules/payment/usecase"
	"main/modules/payment"
	"time"
)

type paymentGrpcHandler struct {
	paymentUsecase usecase.PaymentUsecaseService
	paymentPb.UnimplementedPaymentServiceServer
}

func NewPaymentGrpcHandler(paymentUsecase usecase.PaymentUsecaseService) paymentPb.PaymentServiceServer {
	return &paymentGrpcHandler{
		paymentUsecase: paymentUsecase,
	}
}

func (h *paymentGrpcHandler) CreatePayment(ctx context.Context, req *paymentPb.CreatePaymentRequest) (*paymentPb.PaymentResponse, error) {
    result, err := h.paymentUsecase.CreatePayment(ctx, req.UserId, req.BookingId, req.PaymentMethod, req.FacilityName, req.Amount)
    if err != nil {
        return nil, fmt.Errorf("failed to process payment: %v", err)
    }

	var status string
	switch result.Status {
	case payment.Pending:
		status = "Pending"
	case payment.Completed:
		status = "Completed"
	default:
		status = "Unknown"
	}
	
	
	return &paymentPb.PaymentResponse{
		PaymentId:    result.PaymentID,
		UserId:       result.UserId,
		BookingId:    result.BookingId,
		Amount:       result.Amount,
		Currency:     result.Currency,
		PaymentMethod: result.PaymentMethod,
		Status:       status, // now a string
		FacilityName: result.FacilityName,
		CreatedAt:    result.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    result.UpdatedAt.Format(time.RFC3339),
		QrCodeUrl:    result.QRCodeURL,
	}, nil
	
}


func (h *paymentGrpcHandler) UpdatePaymentStatus(ctx context.Context, req *paymentPb.UpdatePaymentStatusRequest) (*paymentPb.PaymentResponse, error) {
	result, err := h.paymentUsecase.UpdatePayment(ctx, req.PaymentId, req.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to update payment: %v", err)
	}

	var status string
	switch result.Status {
	case payment.Pending:
		status = "Pending"
	case payment.Completed:
		status = "Completed"
	default:
		status = "Unknown"
	}
	

	return &paymentPb.PaymentResponse{
		PaymentId:    result.Id.Hex(),
		UserId:       result.UserID,
		BookingId:    result.BookingID,
		Amount:       result.Amount,
		Currency:     result.Currency,
		PaymentMethod: result.PaymentMethod,
		Status:       status,
		FacilityName: result.FacilityName,
		CreatedAt:    result.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    result.UpdatedAt.Format(time.RFC3339),
		QrCodeUrl:    result.QRCodeURL,
	}, nil
}