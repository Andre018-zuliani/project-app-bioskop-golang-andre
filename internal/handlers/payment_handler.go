package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/andre/project-app-bioskop-golang/internal/middleware"
	"github.com/andre/project-app-bioskop-golang/internal/models"
	"github.com/andre/project-app-bioskop-golang/internal/services"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// PaymentHandler handles payment-related HTTP requests
type PaymentHandler struct {
	paymentService *services.PaymentService
	validator      *validator.Validate
	logger         *zap.Logger
}

// NewPaymentHandler creates a new PaymentHandler
func NewPaymentHandler(paymentService *services.PaymentService, validator *validator.Validate, logger *zap.Logger) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		validator:      validator,
		logger:         logger,
	}
}

// GetPaymentMethods handles getting all payment methods
func (h *PaymentHandler) GetPaymentMethods(w http.ResponseWriter, r *http.Request) {
	// Get payment methods
	methods, err := h.paymentService.GetPaymentMethods(r.Context())
	if err != nil {
		h.logger.Error("failed to get payment methods", zap.Error(err))
		writeError(w, "Failed to get payment methods", http.StatusInternalServerError)
		return
	}

	h.logger.Info("payment methods retrieved successfully")
	writeJSON(w, methods, http.StatusOK)
}

// ProcessPayment handles processing a payment
func (h *PaymentHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(r)
	if err != nil {
		h.logger.Error("failed to get user id from context", zap.Error(err))
		writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request", zap.Error(err))
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		h.logger.Error("validation error", zap.Error(err))
		writeError(w, "Validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Process payment
	response, err := h.paymentService.ProcessPayment(r.Context(), userID, &req)
	if err != nil {
		h.logger.Error("failed to process payment", zap.Error(err), zap.Int("user_id", userID))
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send async payment confirmation using goroutine
	go func() {
		h.logger.Info("Sending payment confirmation notification",
			zap.Int("payment_id", response.ID),
			zap.Int("booking_id", response.BookingID),
			zap.Float64("amount", response.Amount),
		)
		// In production: Send payment receipt via email
		// notificationService.SendPaymentConfirmation(...)
	}()

	h.logger.Info("payment processed successfully", zap.Int("payment_id", response.ID), zap.Int("user_id", userID))
	writeJSON(w, response, http.StatusCreated)
}
