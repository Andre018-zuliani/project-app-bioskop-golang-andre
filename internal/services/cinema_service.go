package services

import (
	"context"
	"fmt"

	"github.com/andre/project-app-bioskop-golang/internal/models"
)

// CinemaService handles cinema-related business logic
type CinemaService struct {
	cinemaRepo CinemaRepository
}

// NewCinemaService creates a new CinemaService
func NewCinemaService(cinemaRepo CinemaRepository) *CinemaService {
	return &CinemaService{cinemaRepo: cinemaRepo}
}

// GetAllCinemas retrieves all cinemas with pagination
func (s *CinemaService) GetAllCinemas(ctx context.Context, page, limit int, filters *models.CinemaFilters) (*models.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	cinemas, total, err := s.cinemaRepo.GetAllCinemas(ctx, page, limit, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get cinemas: %w", err)
	}

	totalPages := (total + limit - 1) / limit

	return &models.PaginatedResponse{
		Data:       cinemas,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// GetCinemaByID retrieves a cinema by ID
func (s *CinemaService) GetCinemaByID(ctx context.Context, id int) (*models.Cinema, error) {
	cinema, err := s.cinemaRepo.GetCinemaByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get cinema: %w", err)
	}
	return cinema, nil
}
