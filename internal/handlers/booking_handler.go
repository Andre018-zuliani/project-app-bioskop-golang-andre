package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/andre/project-app-bioskop-golang/internal/middleware"
	"github.com/andre/project-app-bioskop-golang/internal/models"
	"github.com/andre/project-app-bioskop-golang/internal/services"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// BookingHandler handles booking-related HTTP requests
type BookingHandler struct {
	bookingService *services.BookingService
	validator      *validator.Validate
	logger         *zap.Logger
}

// NewBookingHandler creates a new BookingHandler
func NewBookingHandler(bookingService *services.BookingService, validator *validator.Validate, logger *zap.Logger) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
		validator:      validator,
		logger:         logger,
	}
}

// CreateBooking handles creating a new booking
func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(r)
	if err != nil {
		h.logger.Error("failed to get user id from context", zap.Error(err))
		writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.BookingRequest
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

	// Create booking
	response, err := h.bookingService.CreateBooking(r.Context(), userID, &req)
	if err != nil {
		h.logger.Error("failed to create booking", zap.Error(err), zap.Int("user_id", userID))
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send async notification using goroutine
	go func() {
		h.logger.Info("Sending booking confirmation notification",
			zap.Int("booking_id", response.ID),
			zap.Int("user_id", userID),
		)
		// In production: Send email/SMS notification here
		// notificationService.SendBookingConfirmation(...)
	}()

	h.logger.Info("booking created successfully", zap.Int("booking_id", response.ID), zap.Int("user_id", userID))
	writeJSON(w, response, http.StatusCreated)
}

// GetUserBookings handles getting user bookings
func (h *BookingHandler) GetUserBookings(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(r)
	if err != nil {
		h.logger.Error("failed to get user id from context", zap.Error(err))
		writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get query parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Get user bookings
	response, err := h.bookingService.GetUserBookings(r.Context(), userID, page, limit)
	if err != nil {
		h.logger.Error("failed to get user bookings", zap.Error(err), zap.Int("user_id", userID))
		writeError(w, "Failed to get bookings", http.StatusInternalServerError)
		return
	}

	h.logger.Info("user bookings retrieved successfully", zap.Int("user_id", userID), zap.Int("total", response.Total))
	writeJSON(w, response, http.StatusOK)
}
