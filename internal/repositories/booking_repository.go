package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/andre/project-app-bioskop-golang/internal/models"
	"github.com/jackc/pgx/v5"
)

// BookingRepository handles booking-related database operations
type BookingRepository struct {
	db Database
}

// NewBookingRepository creates a new BookingRepository
func NewBookingRepository(db Database) *BookingRepository {
	return &BookingRepository{db: db}
}

// CreateBooking creates a new booking
func (r *BookingRepository) CreateBooking(ctx context.Context, booking *models.Booking) error {
	query := `INSERT INTO bookings (user_id, cinema_id, seat_id, show_date, show_time, status, total_price, payment_method, payment_status) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, booking_date, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, booking.UserID, booking.CinemaID, booking.SeatID, booking.ShowDate, booking.ShowTime,
		booking.Status, booking.TotalPrice, booking.PaymentMethod, booking.PaymentStatus).
		Scan(&booking.ID, &booking.BookingDate, &booking.CreatedAt, &booking.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create booking: %w", err)
	}
	return nil
}

// GetBookingByID retrieves a booking by ID
func (r *BookingRepository) GetBookingByID(ctx context.Context, id int) (*models.Booking, error) {
	booking := &models.Booking{}
	query := `SELECT id, user_id, cinema_id, seat_id, show_date, show_time, booking_date, status, total_price, 
	payment_method, payment_status, created_at, updated_at FROM bookings WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).
		Scan(&booking.ID, &booking.UserID, &booking.CinemaID, &booking.SeatID, &booking.ShowDate, &booking.ShowTime,
			&booking.BookingDate, &booking.Status, &booking.TotalPrice, &booking.PaymentMethod, &booking.PaymentStatus,
			&booking.CreatedAt, &booking.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	return booking, nil
}

// GetUserBookings retrieves all bookings for a user with pagination
func (r *BookingRepository) GetUserBookings(ctx context.Context, userID int, page, limit int) ([]*models.Booking, int, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) FROM bookings WHERE user_id = $1"
	err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bookings: %w", err)
	}

	// Get paginated data
	query := `SELECT id, user_id, cinema_id, seat_id, show_date, show_time, booking_date, status, total_price, 
	payment_method, payment_status, created_at, updated_at FROM bookings 
	WHERE user_id = $1 ORDER BY booking_date DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user bookings: %w", err)
	}
	defer rows.Close()

	bookings := []*models.Booking{}
	for rows.Next() {
		booking := &models.Booking{}
		err := rows.Scan(&booking.ID, &booking.UserID, &booking.CinemaID, &booking.SeatID, &booking.ShowDate,
			&booking.ShowTime, &booking.BookingDate, &booking.Status, &booking.TotalPrice, &booking.PaymentMethod,
			&booking.PaymentStatus, &booking.CreatedAt, &booking.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, booking)
	}

	return bookings, total, nil
}

// UpdateBookingStatus updates the status of a booking
func (r *BookingRepository) UpdateBookingStatus(ctx context.Context, id int, status string) error {
	query := `UPDATE bookings SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update booking status: %w", err)
	}
	return nil
}

// UpdateBookingPaymentStatus updates the payment status of a booking
func (r *BookingRepository) UpdateBookingPaymentStatus(ctx context.Context, id int, paymentStatus string) error {
	query := `UPDATE bookings SET payment_status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.Exec(ctx, query, paymentStatus, id)
	if err != nil {
		return fmt.Errorf("failed to update booking payment status: %w", err)
	}
	return nil
}

// CheckSeatBooked checks if a seat is already booked for the given date and time
func (r *BookingRepository) CheckSeatBooked(ctx context.Context, seatID int, showDate time.Time, showTime string) (bool, error) {
	query := `SELECT COUNT(*) FROM bookings WHERE seat_id = $1 AND show_date = $2 AND show_time = $3 AND status != 'cancelled'`

	var count int
	err := r.db.QueryRow(ctx, query, seatID, showDate.Format("2006-01-02"), showTime).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check seat booking: %w", err)
	}

	return count > 0, nil
}

// GetBookingWithDetails retrieves booking with cinema and seat details
func (r *BookingRepository) GetBookingWithDetails(ctx context.Context, id int) (*models.Booking, error) {
	booking := &models.Booking{}
	cinema := &models.Cinema{}
	seat := &models.Seat{}

	query := `SELECT b.id, b.user_id, b.cinema_id, b.seat_id, b.show_date, b.show_time, b.booking_date, b.status, b.total_price, 
	b.payment_method, b.payment_status, b.created_at, b.updated_at,
	c.id, c.name, c.location, c.city, c.address, c.total_seats, c.image_url, c.created_at, c.updated_at,
	s.id, s.cinema_id, s.seat_number, s.row_number, s.seat_type, s.price, s.created_at, s.updated_at
	FROM bookings b
	JOIN cinemas c ON b.cinema_id = c.id
	JOIN seats s ON b.seat_id = s.id
	WHERE b.id = $1`

	err := r.db.QueryRow(ctx, query, id).
		Scan(&booking.ID, &booking.UserID, &booking.CinemaID, &booking.SeatID, &booking.ShowDate, &booking.ShowTime,
			&booking.BookingDate, &booking.Status, &booking.TotalPrice, &booking.PaymentMethod, &booking.PaymentStatus,
			&booking.CreatedAt, &booking.UpdatedAt,
			&cinema.ID, &cinema.Name, &cinema.Location, &cinema.City, &cinema.Address, &cinema.TotalSeats, &cinema.ImageURL,
			&cinema.CreatedAt, &cinema.UpdatedAt,
			&seat.ID, &seat.CinemaID, &seat.SeatNumber, &seat.RowNumber, &seat.SeatType, &seat.Price, &seat.CreatedAt, &seat.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get booking with details: %w", err)
	}

	booking.Cinema = cinema
	booking.Seat = seat
	return booking, nil
}
