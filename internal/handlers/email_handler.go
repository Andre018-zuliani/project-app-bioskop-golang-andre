package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/andre/project-app-bioskop-golang/internal/models"
	"github.com/andre/project-app-bioskop-golang/internal/services"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// EmailHandler handles email verification HTTP requests
type EmailHandler struct {
	emailService *services.EmailService
	validator    *validator.Validate
	logger       *zap.Logger
}

// NewEmailHandler creates a new email handler
func NewEmailHandler(emailService *services.EmailService, validator *validator.Validate, logger *zap.Logger) *EmailHandler {
	return &EmailHandler{
		emailService: emailService,
		validator:    validator,
		logger:       logger,
	}
}

// VerifyEmail handles POST /api/verify-email
func (h *EmailHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req models.VerifyOTPRequest

	// Decode request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request", zap.Error(err))
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		h.logger.Error("Validation failed", zap.Error(err))
		writeError(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Verify OTP
	err := h.emailService.VerifyOTP(r.Context(), req.Email, req.OTPCode)
	if err != nil {
		h.logger.Error("OTP verification failed", zap.Error(err), zap.String("email", req.Email))
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Success response
	response := models.OTPResponse{
		Message:    "Email verified successfully",
		IsVerified: true,
	}

	writeJSON(w, response, http.StatusOK)
	h.logger.Info("Email verified", zap.String("email", req.Email))
}

// ResendOTP handles POST /api/resend-otp
func (h *EmailHandler) ResendOTP(w http.ResponseWriter, r *http.Request) {
	var req models.SendOTPRequest

	// Decode request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request", zap.Error(err))
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		h.logger.Error("Validation failed", zap.Error(err))
		writeError(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Resend OTP
	err := h.emailService.ResendOTP(r.Context(), req.Email)
	if err != nil {
		h.logger.Error("Failed to resend OTP", zap.Error(err), zap.String("email", req.Email))
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Success response
	response := models.OTPResponse{
		Message:    "OTP sent to your email",
		IsVerified: false,
	}

	writeJSON(w, response, http.StatusOK)
	h.logger.Info("OTP resent", zap.String("email", req.Email))
}
