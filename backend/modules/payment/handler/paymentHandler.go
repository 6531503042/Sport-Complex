package handler

import (
	"fmt"
	client "main/client/payment"
	"main/config"
	"main/modules/payment"
	"main/modules/payment/usecase"
	"main/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PaymentHttpHandlerService interface {
	CreatePayment(c echo.Context) error
	FindPayment(c echo.Context) error
	FindPaymentsByUser(c echo.Context) error
	HandlePaymentSuccess(c echo.Context) error
	SaveSlip(c echo.Context) error
	UpdateSlipStatus(c echo.Context) error
	GetPendingSlips(c echo.Context) error
}

type paymentHttpHandler struct {
	cfg            *config.Config
	paymentUsecase usecase.PaymentUsecaseService
	paymentClient  *client.PaymentClient
}

// NewPaymentHttpHandler creates a new PaymentHttpHandler
func NewPaymentHttpHandler(cfg *config.Config, paymentUsecase usecase.PaymentUsecaseService, paymentClient *client.PaymentClient) PaymentHttpHandlerService {
	return &paymentHttpHandler{
		cfg:            cfg,
		paymentUsecase: paymentUsecase,
		paymentClient:  paymentClient,
	}
}

// CreatePayment handles the creation of a new payment
func (h *paymentHttpHandler) CreatePayment(c echo.Context) error {
	var req payment.CreatePaymentRequest

	// Bind and validate the request
	if err := c.Bind(&req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "Invalid request format")
	}
	if err := c.Validate(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	// Call the usecase to create the payment
	createdPayment, err := h.paymentUsecase.CreatePayment(c.Request().Context(), req.UserId, req.BookingId, req.PaymentMethod, req.FacilityName , req.Amount)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, "Failed to create payment: "+err.Error())
	}

	// Return the successful response
	return response.SuccessResponse(c, http.StatusCreated, createdPayment)
}

// FindPayment retrieves payment information by ID
func (h *paymentHttpHandler) FindPayment(c echo.Context) error {
	id := c.Param("id")
	payment, err := h.paymentUsecase.FindPayment(c.Request().Context(), id)
	if err != nil {
		return response.ErrResponse(c, http.StatusNotFound, "Payment not found: "+err.Error())
	}
	return response.SuccessResponse(c, http.StatusOK, payment)
}

func (h *paymentHttpHandler) FindPaymentsByUser(c echo.Context) error {
	userId := c.Param("userId") // Get the user ID from the URL parameter

	if userId == "" {
		return response.ErrResponse(c, http.StatusBadRequest, "User ID is required")
	}

	// Call the usecase to retrieve payments
	payments, err := h.paymentUsecase.FindPaymentsByUser(c.Request().Context(), userId)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, "Failed to retrieve payments: "+err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, payments)
}


// HandlePaymentSuccess handles payment success callback
func (h *paymentHttpHandler) HandlePaymentSuccess(c echo.Context) error {
	paymentID := c.Param("payment_id")
	if paymentID == "" {
		return response.ErrResponse(c, http.StatusBadRequest, "Payment ID is required")
	}

	// Update payment status to "COMPLETED"
	err := h.paymentClient.UpdatePaymentStatus(paymentID, "COMPLETED")
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, "Failed to update payment status: "+err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, "Payment status updated to completed")
}

// SaveSlip handles saving a payment slip
func (h *paymentHttpHandler) SaveSlip(c echo.Context) error {
	var slip payment.PaymentSlip
	if err := c.Bind(&slip); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "Invalid input format for slip")
	}

	if err := h.paymentUsecase.SaveSlip(c.Request().Context(), slip); err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, "Failed to save slip: "+err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, "Slip saved successfully")
}

// UpdateSlipStatus handles updating the status of a payment slip
func (h *paymentHttpHandler) UpdateSlipStatus(c echo.Context) error {
	slipId := c.Param("slipId")
	var req payment.UpdateSlipStatusRequest
	if err := c.Bind(&req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "Invalid input format for status update")
	}

	if err := h.paymentUsecase.UpdateSlipStatus(c.Request().Context(), slipId, req.Status); err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, "Failed to update slip status: "+err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, "Slip status updated successfully")
}

// GetPendingSlips handles retrieving all pending payment slips
func (h *paymentHttpHandler) GetPendingSlips(c echo.Context) error {
	slips, err := h.paymentUsecase.GetPendingSlips(c.Request().Context())
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, "Failed to retrieve pending slips: "+err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, slips)
}

// generateQRCodeURL generates a QR code URL for payment
func generateQRCodeURL(payment *payment.PaymentEntity) string {
	baseURL := "https://your-payment-gateway.com/pay"
	return fmt.Sprintf("%s?amount=%.2f&currency=%s&user_id=%s&booking_id=%s",
		baseURL, payment.Amount, payment.Currency, payment.UserID, payment.BookingID)
}
