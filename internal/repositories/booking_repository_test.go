package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/andre/project-app-bioskop-golang/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	pgxmock "github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

type mockDB struct {
	pool pgxmock.PgxPoolIface
}

func (m *mockDB) Close(ctx context.Context) error {
	m.pool.Close()
	return nil
}

func (m *mockDB) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return m.pool.Query(ctx, query, args...)
}

func (m *mockDB) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return m.pool.QueryRow(ctx, query, args...)
}

func (m *mockDB) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return m.pool.Exec(ctx, query, args...)
}

func (m *mockDB) Begin(ctx context.Context) (pgx.Tx, error) {
	return m.pool.Begin(ctx)
}

func TestCreateBooking_Success(t *testing.T) {
	pool, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer pool.Close()

	repo := NewBookingRepository(&mockDB{pool: pool})

	booking := &models.Booking{
		UserID:        1,
		CinemaID:      1,
		SeatID:        1,
		ShowDate:      time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
		ShowTime:      "19:00",
		Status:        "pending",
		TotalPrice:    50000,
		PaymentMethod: "cash",
		PaymentStatus: "pending",
	}

	pool.ExpectQuery("INSERT INTO bookings").
		WithArgs(booking.UserID, booking.CinemaID, booking.SeatID, booking.ShowDate, booking.ShowTime, booking.Status, booking.TotalPrice, booking.PaymentMethod, booking.PaymentStatus).
		WillReturnRows(pgxmock.NewRows([]string{"id", "booking_date", "created_at", "updated_at"}).AddRow(1, time.Now(), time.Now(), time.Now()))

	err = repo.CreateBooking(context.Background(), booking)

	assert.NoError(t, err)
	assert.Equal(t, 1, booking.ID)
	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestCheckSeatBooked(t *testing.T) {
	pool, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer pool.Close()

	repo := NewBookingRepository(&mockDB{pool: pool})
	showDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

	pool.ExpectQuery(`SELECT COUNT\(\*\) FROM bookings`).
		WithArgs(1, showDate.Format("2006-01-02"), "19:00").
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(1))

	booked, err := repo.CheckSeatBooked(context.Background(), 1, showDate, "19:00")

	assert.NoError(t, err)
	assert.True(t, booked)
	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestGetBookingByID_NoRows(t *testing.T) {
	pool, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer pool.Close()

	repo := NewBookingRepository(&mockDB{pool: pool})

	pool.ExpectQuery("SELECT id, user_id").WithArgs(99).WillReturnError(pgx.ErrNoRows)

	booking, err := repo.GetBookingByID(context.Background(), 99)

	assert.NoError(t, err)
	assert.Nil(t, booking)
	assert.NoError(t, pool.ExpectationsWereMet())
}
