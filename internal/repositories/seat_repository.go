package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/andre/project-app-bioskop-golang/internal/models"
	"github.com/jackc/pgx/v5"
)

// SeatRepository handles seat-related database operations
type SeatRepository struct {
	db Database
}

// NewSeatRepository creates a new SeatRepository
func NewSeatRepository(db Database) *SeatRepository {
	return &SeatRepository{db: db}
}

// GetSeatsByCinema retrieves all seats for a cinema
func (r *SeatRepository) GetSeatsByCinema(ctx context.Context, cinemaID int) ([]*models.Seat, error) {
	query := `SELECT id, cinema_id, seat_number, row_number, seat_type, price, created_at, updated_at 
	FROM seats WHERE cinema_id = $1 ORDER BY row_number, seat_number`

	rows, err := r.db.Query(ctx, query, cinemaID)
	if err != nil {
		return nil, fmt.Errorf("failed to get seats: %w", err)
	}
	defer rows.Close()

	seats := []*models.Seat{}
	for rows.Next() {
		seat := &models.Seat{}
		err := rows.Scan(&seat.ID, &seat.CinemaID, &seat.SeatNumber, &seat.RowNumber, &seat.SeatType, &seat.Price, &seat.CreatedAt, &seat.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan seat: %w", err)
		}
		seats = append(seats, seat)
	}

	return seats, nil
}

// GetSeatByID retrieves a seat by ID
func (r *SeatRepository) GetSeatByID(ctx context.Context, id int) (*models.Seat, error) {
	seat := &models.Seat{}
	query := `SELECT id, cinema_id, seat_number, row_number, seat_type, price, created_at, updated_at 
	FROM seats WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).
		Scan(&seat.ID, &seat.CinemaID, &seat.SeatNumber, &seat.RowNumber, &seat.SeatType, &seat.Price, &seat.CreatedAt, &seat.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get seat: %w", err)
	}

	return seat, nil
}

// CreateSeat creates a new seat (for seeding)
func (r *SeatRepository) CreateSeat(ctx context.Context, seat *models.Seat) error {
	query := `INSERT INTO seats (cinema_id, seat_number, row_number, seat_type, price) 
	VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, seat.CinemaID, seat.SeatNumber, seat.RowNumber, seat.SeatType, seat.Price).
		Scan(&seat.ID, &seat.CreatedAt, &seat.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create seat: %w", err)
	}
	return nil
}

// GetSeatAvailability retrieves seat availability for a specific date and time
func (r *SeatRepository) GetSeatAvailability(ctx context.Context, cinemaID int, date time.Time, timeStr string) ([]*models.SeatAvailability, error) {
	query := `SELECT sa.id, sa.cinema_id, sa.seat_id, sa.show_date, sa.show_time, sa.is_available, sa.created_at, sa.updated_at,
	s.id, s.cinema_id, s.seat_number, s.row_number, s.seat_type, s.price, s.created_at, s.updated_at
	FROM seat_availability sa
	JOIN seats s ON sa.seat_id = s.id
	WHERE sa.cinema_id = $1 AND sa.show_date = $2 AND sa.show_time = $3
	ORDER BY s.row_number, s.seat_number`

	rows, err := r.db.Query(ctx, query, cinemaID, date.Format("2006-01-02"), timeStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get seat availability: %w", err)
	}
	defer rows.Close()

	availabilities := []*models.SeatAvailability{}
	for rows.Next() {
		sa := &models.SeatAvailability{}
		seat := &models.Seat{}
		err := rows.Scan(
			&sa.ID, &sa.CinemaID, &sa.SeatID, &sa.ShowDate, &sa.ShowTime, &sa.IsAvailable, &sa.CreatedAt, &sa.UpdatedAt,
			&seat.ID, &seat.CinemaID, &seat.SeatNumber, &seat.RowNumber, &seat.SeatType, &seat.Price, &seat.CreatedAt, &seat.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan seat availability: %w", err)
		}
		sa.Seat = seat
		availabilities = append(availabilities, sa)
	}

	return availabilities, nil
}

// CreateSeatAvailability creates seat availability records (for seeding)
func (r *SeatRepository) CreateSeatAvailability(ctx context.Context, cinemaID int, date time.Time, timeStr string) error {
	// Get all seats for the cinema
	seats, err := r.GetSeatsBycinema(ctx, cinemaID)
	if err != nil {
		return err
	}

	for _, seat := range seats {
		query := `INSERT INTO seat_availability (cinema_id, seat_id, show_date, show_time, is_available) 
		VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`

		_, err := r.db.Exec(ctx, query, cinemaID, seat.ID, date.Format("2006-01-02"), timeStr, true)
		if err != nil {
			return fmt.Errorf("failed to create seat availability: %w", err)
		}
	}

	return nil
}

// GetSeatsBycinema is a helper to get seats by cinema
func (r *SeatRepository) GetSeatsBycinema(ctx context.Context, cinemaID int) ([]*models.Seat, error) {
	return r.GetSeatsByCinema(ctx, cinemaID)
}

// UpdateSeatAvailability updates the availability of a seat
func (r *SeatRepository) UpdateSeatAvailability(ctx context.Context, seatID int, date time.Time, timeStr string, isAvailable bool) error {
	query := `UPDATE seat_availability SET is_available = $1, updated_at = CURRENT_TIMESTAMP 
	WHERE seat_id = $2 AND show_date = $3 AND show_time = $4`

	_, err := r.db.Exec(ctx, query, isAvailable, seatID, date.Format("2006-01-02"), timeStr)
	if err != nil {
		return fmt.Errorf("failed to update seat availability: %w", err)
	}

	return nil
}
