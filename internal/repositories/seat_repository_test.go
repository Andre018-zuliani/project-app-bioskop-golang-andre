package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/andre/project-app-bioskop-golang/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestSeatRepository_GetSeatsByCinema_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewSeatRepository(&mockDB{pool: mock})

	now := time.Now()
	rows := pgxmock.NewRows([]string{"id", "cinema_id", "seat_number", "row_number", "seat_type", "price", "created_at", "updated_at"}).
		AddRow(1, 1, "A1", 1, "regular", 50000.0, now, now).
		AddRow(2, 1, "A2", 1, "regular", 50000.0, now, now).
		AddRow(3, 1, "B1", 2, "vip", 100000.0, now, now)

	mock.ExpectQuery("SELECT id, cinema_id, seat_number").
		WithArgs(1).
		WillReturnRows(rows)

	// Execute
	seats, err := repo.GetSeatsByCinema(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, seats, 3)
	assert.Equal(t, "A1", seats[0].SeatNumber)
	assert.Equal(t, "B1", seats[2].SeatNumber)
	assert.Equal(t, "vip", seats[2].SeatType)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSeatRepository_GetSeatByID_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewSeatRepository(&mockDB{pool: mock})

	now := time.Now()
	rows := pgxmock.NewRows([]string{"id", "cinema_id", "seat_number", "row_number", "seat_type", "price", "created_at", "updated_at"}).
		AddRow(1, 1, "A1", 1, "regular", 50000.0, now, now)

	mock.ExpectQuery("SELECT id, cinema_id, seat_number").
		WithArgs(1).
		WillReturnRows(rows)

	// Execute
	seat, err := repo.GetSeatByID(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, seat)
	assert.Equal(t, 1, seat.ID)
	assert.Equal(t, "A1", seat.SeatNumber)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSeatRepository_GetSeatByID_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewSeatRepository(&mockDB{pool: mock})

	mock.ExpectQuery("SELECT id, cinema_id, seat_number").
		WithArgs(999).
		WillReturnError(pgx.ErrNoRows)

	// Execute
	seat, err := repo.GetSeatByID(context.Background(), 999)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, seat)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSeatRepository_CreateSeat_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewSeatRepository(&mockDB{pool: mock})

	now := time.Now()
	rows := pgxmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(1, now, now)

	mock.ExpectQuery("INSERT INTO seats").
		WithArgs(1, "A1", 1, "regular", 50000.0).
		WillReturnRows(rows)

	// Execute
	seat := &models.Seat{
		CinemaID:   1,
		SeatNumber: "A1",
		RowNumber:  1,
		SeatType:   "regular",
		Price:      50000.0,
	}
	err = repo.CreateSeat(context.Background(), seat)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, seat.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSeatRepository_GetSeatAvailability_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewSeatRepository(&mockDB{pool: mock})

	now := time.Now()
	showDate := time.Now().AddDate(0, 0, 1) // Tomorrow
	rows := pgxmock.NewRows([]string{
		"sa_id", "sa_cinema_id", "sa_seat_id", "sa_show_date", "sa_show_time", "sa_is_available", "sa_created_at", "sa_updated_at",
		"s_id", "s_cinema_id", "s_seat_number", "s_row_number", "s_seat_type", "s_price", "s_created_at", "s_updated_at",
	}).
		AddRow(1, 1, 1, showDate, "19:00", true, now, now,
			1, 1, "A1", 1, "regular", 50000.0, now, now).
		AddRow(2, 1, 2, showDate, "19:00", false, now, now,
			2, 1, "A2", 1, "regular", 50000.0, now, now)

	mock.ExpectQuery("SELECT sa.id, sa.cinema_id").
		WithArgs(1, showDate.Format("2006-01-02"), "19:00").
		WillReturnRows(rows)

	// Execute
	availabilities, err := repo.GetSeatAvailability(context.Background(), 1, showDate, "19:00")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, availabilities, 2)
	assert.True(t, availabilities[0].IsAvailable)
	assert.False(t, availabilities[1].IsAvailable)
	assert.Equal(t, "A1", availabilities[0].Seat.SeatNumber)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSeatRepository_UpdateSeatAvailability_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewSeatRepository(&mockDB{pool: mock})

	showDate := time.Now().AddDate(0, 0, 1)
	mock.ExpectExec("UPDATE seat_availability SET is_available").
		WithArgs(false, 1, showDate.Format("2006-01-02"), "19:00").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	// Execute
	err = repo.UpdateSeatAvailability(context.Background(), 1, showDate, "19:00", false)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSeatRepository_CreateSeatAvailability_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewSeatRepository(&mockDB{pool: mock})

	now := time.Now()
	showDate := time.Now().AddDate(0, 0, 1)

	// Mock GetSeatsByCinema query first
	seatRows := pgxmock.NewRows([]string{"id", "cinema_id", "seat_number", "row_number", "seat_type", "price", "created_at", "updated_at"}).
		AddRow(1, 1, "A1", 1, "regular", 50000.0, now, now).
		AddRow(2, 1, "A2", 1, "regular", 50000.0, now, now)

	mock.ExpectQuery("SELECT id, cinema_id, seat_number").
		WithArgs(1).
		WillReturnRows(seatRows)

	// Mock INSERT for each seat
	mock.ExpectExec("INSERT INTO seat_availability").
		WithArgs(1, 1, showDate.Format("2006-01-02"), "19:00", true).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	mock.ExpectExec("INSERT INTO seat_availability").
		WithArgs(1, 2, showDate.Format("2006-01-02"), "19:00", true).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	// Execute
	err = repo.CreateSeatAvailability(context.Background(), 1, showDate, "19:00")

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
