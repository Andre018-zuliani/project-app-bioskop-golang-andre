package handlers

import (
	"net/http"
	"strconv"

	"github.com/andre/project-app-bioskop-golang/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// SeatHandler handles seat-related HTTP requests
type SeatHandler struct {
	seatService *services.SeatService
	validator   *validator.Validate
	logger      *zap.Logger
}

// NewSeatHandler creates a new SeatHandler
func NewSeatHandler(seatService *services.SeatService, validator *validator.Validate, logger *zap.Logger) *SeatHandler {
	return &SeatHandler{
		seatService: seatService,
		validator:   validator,
		logger:      logger,
	}
}

// GetSeatAvailability handles getting seat availability
func (h *SeatHandler) GetSeatAvailability(w http.ResponseWriter, r *http.Request) {
	cinemaID := chi.URLParam(r, "cinemaId")
	id, err := strconv.Atoi(cinemaID)
	if err != nil {
		writeError(w, "Invalid cinema ID", http.StatusBadRequest)
		return
	}

	// Get query parameters
	date := r.URL.Query().Get("date")
	time := r.URL.Query().Get("time")

	if date == "" || time == "" {
		writeError(w, "Missing date or time parameter", http.StatusBadRequest)
		return
	}

	// Get seat availability
	response, err := h.seatService.GetSeatAvailability(r.Context(), id, date, time)
	if err != nil {
		h.logger.Error("failed to get seat availability", zap.Error(err), zap.Int("cinema_id", id))
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Info("seat availability retrieved successfully", zap.Int("cinema_id", id), zap.String("date", date), zap.String("time", time))
	writeJSON(w, response, http.StatusOK)
}
