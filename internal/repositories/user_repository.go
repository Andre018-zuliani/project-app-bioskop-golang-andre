package repositories

import (
	"context"
	"fmt"

	"github.com/andre/project-app-bioskop-golang/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Database interface for database operations
type Database interface {
	Close(ctx context.Context) error
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

// UserRepository handles user-related database operations
type UserRepository struct {
	db Database
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db Database) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (username, email, password, is_verified) 
	VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, user.Username, user.Email, user.Password, false).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetUserByUsername retrieves a user by username
func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, is_verified, created_at, updated_at 
	FROM users WHERE username = $1`

	err := r.db.QueryRow(ctx, query, username).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return user, nil
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, is_verified, created_at, updated_at 
	FROM users WHERE email = $1`

	err := r.db.QueryRow(ctx, query, email).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, is_verified, created_at, updated_at 
	FROM users WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return user, nil
}

// CreateSession creates a new user session
func (r *UserRepository) CreateSession(ctx context.Context, session *models.UserSession) error {
	query := `INSERT INTO user_sessions (user_id, token, expires_at) 
	VALUES ($1, $2, $3) RETURNING id, created_at`

	err := r.db.QueryRow(ctx, query, session.UserID, session.Token, session.ExpiresAt).
		Scan(&session.ID, &session.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

// GetSessionByToken retrieves a session by token
func (r *UserRepository) GetSessionByToken(ctx context.Context, token string) (*models.UserSession, error) {
	session := &models.UserSession{}
	query := `SELECT id, user_id, token, created_at, expires_at 
	FROM user_sessions WHERE token = $1`

	err := r.db.QueryRow(ctx, query, token).
		Scan(&session.ID, &session.UserID, &session.Token, &session.CreatedAt, &session.ExpiresAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get session by token: %w", err)
	}
	return session, nil
}

// DeleteSession deletes a session
func (r *UserRepository) DeleteSession(ctx context.Context, token string) error {
	query := `DELETE FROM user_sessions WHERE token = $1`
	_, err := r.db.Exec(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}
