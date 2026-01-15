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

func TestCinemaRepository_GetAllCinemas_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewCinemaRepository(&mockDB{pool: mock})

	// Mock count query
	countRows := pgxmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery("SELECT COUNT").
		WillReturnRows(countRows)

	// Mock data query
	now := time.Now()
	rows := pgxmock.NewRows([]string{"id", "name", "location", "city", "address", "total_seats", "image_url", "created_at", "updated_at"}).
		AddRow(1, "Cinema XXI Plaza", "Jakarta", "Jakarta", "Jl. Sudirman No. 1", 100, "cinema1.jpg", now, now).
		AddRow(2, "Cinema CGV Bandung", "Bandung", "Bandung", "Jl. Asia Afrika No. 10", 150, "cinema2.jpg", now, now)

	mock.ExpectQuery("SELECT id, name, location").
		WithArgs(10, 0).
		WillReturnRows(rows)

	// Execute
	filters := &models.CinemaFilters{}
	cinemas, total, err := repo.GetAllCinemas(context.Background(), 1, 10, filters)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, cinemas, 2)
	assert.Equal(t, "Cinema XXI Plaza", cinemas[0].Name)
	assert.Equal(t, "Cinema CGV Bandung", cinemas[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCinemaRepository_GetAllCinemas_WithFilters(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewCinemaRepository(&mockDB{pool: mock})

	// Mock count query with filter
	countRows := pgxmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery("SELECT COUNT").
		WithArgs("%Jakarta%").
		WillReturnRows(countRows)

	// Mock data query with filter
	now := time.Now()
	rows := pgxmock.NewRows([]string{"id", "name", "location", "city", "address", "total_seats", "image_url", "created_at", "updated_at"}).
		AddRow(1, "Cinema XXI Plaza", "Jakarta", "Jakarta", "Jl. Sudirman No. 1", 100, "cinema1.jpg", now, now)

	mock.ExpectQuery("SELECT id, name, location").
		WithArgs("%Jakarta%", 10, 0).
		WillReturnRows(rows)

	// Execute
	filters := &models.CinemaFilters{City: "Jakarta"}
	cinemas, total, err := repo.GetAllCinemas(context.Background(), 1, 10, filters)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, cinemas, 1)
	assert.Equal(t, "Cinema XXI Plaza", cinemas[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCinemaRepository_GetCinemaByID_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewCinemaRepository(&mockDB{pool: mock})

	now := time.Now()
	rows := pgxmock.NewRows([]string{"id", "name", "location", "city", "address", "total_seats", "image_url", "created_at", "updated_at"}).
		AddRow(1, "Cinema XXI Plaza", "Jakarta", "Jakarta", "Jl. Sudirman No. 1", 100, "cinema1.jpg", now, now)

	mock.ExpectQuery("SELECT id, name, location").
		WithArgs(1).
		WillReturnRows(rows)

	// Execute
	cinema, err := repo.GetCinemaByID(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, cinema)
	assert.Equal(t, 1, cinema.ID)
	assert.Equal(t, "Cinema XXI Plaza", cinema.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCinemaRepository_GetCinemaByID_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewCinemaRepository(&mockDB{pool: mock})

	mock.ExpectQuery("SELECT id, name, location").
		WithArgs(999).
		WillReturnError(pgx.ErrNoRows)

	// Execute
	cinema, err := repo.GetCinemaByID(context.Background(), 999)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, cinema)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCinemaRepository_CreateCinema_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewCinemaRepository(&mockDB{pool: mock})

	now := time.Now()
	rows := pgxmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(1, now, now)

	mock.ExpectQuery("INSERT INTO cinemas").
		WithArgs("Cinema XXI Plaza", "Jakarta", "Jakarta", "Jl. Sudirman No. 1", 100, "cinema1.jpg").
		WillReturnRows(rows)

	// Execute
	cinema := &models.Cinema{
		Name:       "Cinema XXI Plaza",
		Location:   "Jakarta",
		City:       "Jakarta",
		Address:    "Jl. Sudirman No. 1",
		TotalSeats: 100,
		ImageURL:   "cinema1.jpg",
	}
	err = repo.CreateCinema(context.Background(), cinema)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, cinema.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}
