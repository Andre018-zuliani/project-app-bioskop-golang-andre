package handlers

import (
	"net/http"
	"strconv"

	"github.com/andre/project-app-bioskop-golang/internal/models"
	"github.com/andre/project-app-bioskop-golang/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// CinemaHandler handles cinema-related HTTP requests
type CinemaHandler struct {
	cinemaService *services.CinemaService
	validator     *validator.Validate
	logger        *zap.Logger
}

// NewCinemaHandler creates a new CinemaHandler
func NewCinemaHandler(cinemaService *services.CinemaService, validator *validator.Validate, logger *zap.Logger) *CinemaHandler {
	return &CinemaHandler{
		cinemaService: cinemaService,
		validator:     validator,
		logger:        logger,
	}
}

// GetAllCinemas handles getting all cinemas
func (h *CinemaHandler) GetAllCinemas(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	city := r.URL.Query().Get("city")
	name := r.URL.Query().Get("name")

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

	filters := &models.CinemaFilters{
		Page:  page,
		Limit: limit,
		City:  city,
		Name:  name,
	}

	// Get cinemas
	response, err := h.cinemaService.GetAllCinemas(r.Context(), page, limit, filters)
	if err != nil {
		h.logger.Error("failed to get cinemas", zap.Error(err))
		writeError(w, "Failed to get cinemas", http.StatusInternalServerError)
		return
	}

	h.logger.Info("cinemas retrieved successfully", zap.Int("total", response.Total))
	writeJSON(w, response, http.StatusOK)
}

// GetCinemaByID handles getting a cinema by ID
func (h *CinemaHandler) GetCinemaByID(w http.ResponseWriter, r *http.Request) {
	cinemaID := chi.URLParam(r, "cinemaId")
	id, err := strconv.Atoi(cinemaID)
	if err != nil {
		writeError(w, "Invalid cinema ID", http.StatusBadRequest)
		return
	}

	// Get cinema
	cinema, err := h.cinemaService.GetCinemaByID(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to get cinema", zap.Error(err), zap.Int("cinema_id", id))
		writeError(w, "Failed to get cinema", http.StatusInternalServerError)
		return
	}

	if cinema == nil {
		writeError(w, "Cinema not found", http.StatusNotFound)
		return
	}

	h.logger.Info("cinema retrieved successfully", zap.Int("cinema_id", id))
	writeJSON(w, cinema, http.StatusOK)
}
