package handler

import (
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
}

type paymentHttpHandler struct {
	cfg            *config.Config
	paymentUsecase usecase.PaymentUsecaseService
}

// NewPaymentHttpHandler creates a new PaymentHttpHandler
func NewPaymentHttpHandler(cfg *config.Config, paymentUsecase usecase.PaymentUsecaseService) PaymentHttpHandlerService {
	return &paymentHttpHandler{
		cfg:            cfg,
		paymentUsecase: paymentUsecase,
	}
}

// CreatePayment handles the creation of a new payment
func (h *paymentHttpHandler) CreatePayment(c echo.Context) error {
	var req payment.CreatePaymentRequest

	// Bind and validate the request
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Call the usecase to create the payment
	createdPayment, err := h.paymentUsecase.CreatePayment(c.Request().Context(), req.UserId, req.BookingId, req.Amount)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Return the successful response
	return c.JSON(http.StatusCreated, createdPayment)
}

// GetPayment retrieves payment information by ID
func (h *paymentHttpHandler) FindPayment(c echo.Context) error {
	Id := c.Param("id")
	payment, err := h.paymentUsecase.FindPayment(c.Request().Context(), Id)
	if err != nil {
		return response.ErrResponse(c, http.StatusNotFound, err.Error())
	}
	return response.SuccessResponse(c, http.StatusOK, payment)
}
