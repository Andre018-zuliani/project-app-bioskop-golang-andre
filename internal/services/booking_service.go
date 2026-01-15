package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/andre/project-app-bioskop-golang/internal/models"
)

// BookingService handles booking-related business logic
type BookingService struct {
	bookingRepo BookingRepository
	seatRepo    SeatRepository
	cinemaRepo  CinemaRepository
}

// NewBookingService creates a new BookingService
func NewBookingService(bookingRepo BookingRepository, seatRepo SeatRepository, cinemaRepo CinemaRepository) *BookingService {
	return &BookingService{
		bookingRepo: bookingRepo,
		seatRepo:    seatRepo,
		cinemaRepo:  cinemaRepo,
	}
}

// CreateBooking creates a new booking
func (s *BookingService) CreateBooking(ctx context.Context, userID int, req *models.BookingRequest) (*models.BookingResponse, error) {
	// Parse date
	showDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	// Check if seat exists and get its price
	seat, err := s.seatRepo.GetSeatByID(ctx, req.SeatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get seat: %w", err)
	}
	if seat == nil {
		return nil, errors.New("seat not found")
	}

	// Check if cinema exists
	cinema, err := s.cinemaRepo.GetCinemaByID(ctx, req.CinemaID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cinema: %w", err)
	}
	if cinema == nil {
		return nil, errors.New("cinema not found")
	}

	// Check if seat belongs to the cinema
	if seat.CinemaID != req.CinemaID {
		return nil, errors.New("seat does not belong to this cinema")
	}

	// Check if seat is already booked
	isBooked, err := s.bookingRepo.CheckSeatBooked(ctx, req.SeatID, showDate, req.Time)
	if err != nil {
		return nil, fmt.Errorf("failed to check seat booking: %w", err)
	}
	if isBooked {
		return nil, errors.New("seat is already booked for this date and time")
	}

	// Create booking
	booking := &models.Booking{
		UserID:        userID,
		CinemaID:      req.CinemaID,
		SeatID:        req.SeatID,
		ShowDate:      showDate,
		ShowTime:      req.Time,
		Status:        "pending",
		TotalPrice:    seat.Price,
		PaymentMethod: req.PaymentMethod,
		PaymentStatus: "pending",
	}

	err = s.bookingRepo.CreateBooking(ctx, booking)
	if err != nil {
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	// Update seat availability
	err = s.seatRepo.UpdateSeatAvailability(ctx, req.SeatID, showDate, req.Time, false)
	if err != nil {
		return nil, fmt.Errorf("failed to update seat availability: %w", err)
	}

	response := &models.BookingResponse{
		ID:            booking.ID,
		CinemaID:      booking.CinemaID,
		SeatID:        booking.SeatID,
		ShowDate:      booking.ShowDate,
		ShowTime:      booking.ShowTime,
		TotalPrice:    booking.TotalPrice,
		PaymentMethod: booking.PaymentMethod,
		Status:        booking.Status,
		PaymentStatus: booking.PaymentStatus,
		CreatedAt:     booking.CreatedAt,
	}

	return response, nil
}

// GetUserBookings retrieves bookings for a user
func (s *BookingService) GetUserBookings(ctx context.Context, userID int, page, limit int) (*models.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	bookings, total, err := s.bookingRepo.GetUserBookings(ctx, userID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user bookings: %w", err)
	}

	// Enrich bookings with cinema and seat details
	enrichedBookings := make([]*models.Booking, 0, len(bookings))
	for _, booking := range bookings {
		fullBooking, err := s.bookingRepo.GetBookingWithDetails(ctx, booking.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get booking details: %w", err)
		}
		if fullBooking != nil {
			enrichedBookings = append(enrichedBookings, fullBooking)
		}
	}

	totalPages := (total + limit - 1) / limit

	return &models.PaginatedResponse{
		Data:       enrichedBookings,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// GetBookingByID retrieves a booking by ID
func (s *BookingService) GetBookingByID(ctx context.Context, id int) (*models.Booking, error) {
	booking, err := s.bookingRepo.GetBookingWithDetails(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}
	return booking, nil
}

// UpdateBookingStatus updates the status of a booking
func (s *BookingService) UpdateBookingStatus(ctx context.Context, id int, status string) error {
	err := s.bookingRepo.UpdateBookingStatus(ctx, id, status)
	if err != nil {
		return fmt.Errorf("failed to update booking status: %w", err)
	}
	return nil
}
