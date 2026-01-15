package services

import (
	"context"
	"fmt"
	"time"

	"github.com/andre/project-app-bioskop-golang/internal/models"
)

// SeatService handles seat-related business logic
type SeatService struct {
	seatRepo SeatRepository
}

// NewSeatService creates a new SeatService
func NewSeatService(seatRepo SeatRepository) *SeatService {
	return &SeatService{seatRepo: seatRepo}
}

// GetSeatAvailability retrieves seat availability for a specific cinema, date, and time
func (s *SeatService) GetSeatAvailability(ctx context.Context, cinemaID int, dateStr, timeStr string) (*models.SeatAvailabilityResponse, error) {
	// Parse date
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	// Get seat availability from repository
	availabilities, err := s.seatRepo.GetSeatAvailability(ctx, cinemaID, date, timeStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get seat availability: %w", err)
	}

	// Separate available and unavailable seats
	availableSeats := []*models.SeatAvailability{}
	unavailableSeats := []*models.SeatAvailability{}

	for _, sa := range availabilities {
		if sa.IsAvailable {
			availableSeats = append(availableSeats, sa)
		} else {
			unavailableSeats = append(unavailableSeats, sa)
		}
	}

	response := &models.SeatAvailabilityResponse{
		CinemaID:         cinemaID,
		Date:             dateStr,
		Time:             timeStr,
		AvailableSeats:   availableSeats,
		UnavailableSeats: unavailableSeats,
		TotalAvailable:   len(availableSeats),
		TotalUnavailable: len(unavailableSeats),
	}

	return response, nil
}

// GetSeatByID retrieves a seat by ID
func (s *SeatService) GetSeatByID(ctx context.Context, id int) (*models.Seat, error) {
	seat, err := s.seatRepo.GetSeatByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get seat: %w", err)
	}
	return seat, nil
}
