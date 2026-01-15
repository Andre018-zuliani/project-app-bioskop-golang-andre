package repositories

import (
	"context"
	"fmt"

	"github.com/andre/project-app-bioskop-golang/internal/models"
	"github.com/jackc/pgx/v5"
)

// CinemaRepository handles cinema-related database operations
type CinemaRepository struct {
	db Database
}

// NewCinemaRepository creates a new CinemaRepository
func NewCinemaRepository(db Database) *CinemaRepository {
	return &CinemaRepository{db: db}
}

// GetAllCinemas retrieves all cinemas with pagination
func (r *CinemaRepository) GetAllCinemas(ctx context.Context, page, limit int, filters *models.CinemaFilters) ([]*models.Cinema, int, error) {
	offset := (page - 1) * limit

	// Build WHERE clause
	whereClause := ""
	args := []interface{}{}
	argIndex := 1

	if filters.City != "" {
		whereClause += fmt.Sprintf(" WHERE city ILIKE $%d", argIndex)
		args = append(args, "%"+filters.City+"%")
		argIndex++
	}

	if filters.Name != "" {
		if whereClause == "" {
			whereClause += fmt.Sprintf(" WHERE name ILIKE $%d", argIndex)
		} else {
			whereClause += fmt.Sprintf(" AND name ILIKE $%d", argIndex)
		}
		args = append(args, "%"+filters.Name+"%")
		argIndex++
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM cinemas" + whereClause
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count cinemas: %w", err)
	}

	// Get paginated data
	query := fmt.Sprintf("SELECT id, name, location, city, address, total_seats, image_url, created_at, updated_at "+
		"FROM cinemas%s ORDER BY name ASC LIMIT $%d OFFSET $%d", whereClause, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get cinemas: %w", err)
	}
	defer rows.Close()

	cinemas := []*models.Cinema{}
	for rows.Next() {
		cinema := &models.Cinema{}
		err := rows.Scan(&cinema.ID, &cinema.Name, &cinema.Location, &cinema.City, &cinema.Address, &cinema.TotalSeats, &cinema.ImageURL, &cinema.CreatedAt, &cinema.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan cinema: %w", err)
		}
		cinemas = append(cinemas, cinema)
	}

	return cinemas, total, nil
}

// GetCinemaByID retrieves a cinema by ID
func (r *CinemaRepository) GetCinemaByID(ctx context.Context, id int) (*models.Cinema, error) {
	cinema := &models.Cinema{}
	query := `SELECT id, name, location, city, address, total_seats, image_url, created_at, updated_at 
	FROM cinemas WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).
		Scan(&cinema.ID, &cinema.Name, &cinema.Location, &cinema.City, &cinema.Address, &cinema.TotalSeats, &cinema.ImageURL, &cinema.CreatedAt, &cinema.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get cinema: %w", err)
	}

	return cinema, nil
}

// CreateCinema creates a new cinema (for admin/seeding)
func (r *CinemaRepository) CreateCinema(ctx context.Context, cinema *models.Cinema) error {
	query := `INSERT INTO cinemas (name, location, city, address, total_seats, image_url) 
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, cinema.Name, cinema.Location, cinema.City, cinema.Address, cinema.TotalSeats, cinema.ImageURL).
		Scan(&cinema.ID, &cinema.CreatedAt, &cinema.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create cinema: %w", err)
	}
	return nil
}
