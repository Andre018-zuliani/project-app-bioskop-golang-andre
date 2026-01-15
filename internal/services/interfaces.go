package services

import (
	"context"
	"time"

	"github.com/andre/project-app-bioskop-golang/internal/models"
)

// BookingRepository defines the data access behavior needed by booking-related services.
type BookingRepository interface {
	CreateBooking(ctx context.Context, booking *models.Booking) error
	GetBookingByID(ctx context.Context, id int) (*models.Booking, error)
	GetBookingWithDetails(ctx context.Context, id int) (*models.Booking, error)
	GetUserBookings(ctx context.Context, userID, page, limit int) ([]*models.Booking, int, error)
	UpdateBookingStatus(ctx context.Context, id int, status string) error
	UpdateBookingPaymentStatus(ctx context.Context, id int, paymentStatus string) error
	CheckSeatBooked(ctx context.Context, seatID int, showDate time.Time, showTime string) (bool, error)
}

// SeatRepository describes seat persistence behaviors.
type SeatRepository interface {
	GetSeatAvailability(ctx context.Context, cinemaID int, date time.Time, timeStr string) ([]*models.SeatAvailability, error)
	GetSeatByID(ctx context.Context, id int) (*models.Seat, error)
	UpdateSeatAvailability(ctx context.Context, seatID int, date time.Time, timeStr string, isAvailable bool) error
}

// CinemaRepository defines the storage behavior for cinemas used by services.
type CinemaRepository interface {
	GetAllCinemas(ctx context.Context, page, limit int, filters *models.CinemaFilters) ([]*models.Cinema, int, error)
	GetCinemaByID(ctx context.Context, id int) (*models.Cinema, error)
}

// PaymentRepository describes payment persistence behaviors.
type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *models.Payment) error
	GetPaymentByID(ctx context.Context, id int) (*models.Payment, error)
	GetPaymentByBookingID(ctx context.Context, bookingID int) (*models.Payment, error)
	UpdatePaymentStatus(ctx context.Context, id int, status string) error
	GetPaymentMethods(ctx context.Context) ([]*models.PaymentMethod, error)
	GetPaymentMethodByName(ctx context.Context, name string) (*models.PaymentMethod, error)
}
