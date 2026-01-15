package repositories

import (
	"context"

	"github.com/andre/project-app-bioskop-golang/internal/models"
	"github.com/jackc/pgx/v5"
)

// EmailVerificationRepository handles email verification data operations
type EmailVerificationRepository struct {
	conn *pgx.Conn
}

// NewEmailVerificationRepository creates a new repository
func NewEmailVerificationRepository(conn *pgx.Conn) *EmailVerificationRepository {
	return &EmailVerificationRepository{conn: conn}
}

// Create saves a new email verification record
func (r *EmailVerificationRepository) Create(ctx context.Context, verification *models.EmailVerification) error {
	query := `
		INSERT INTO email_verifications (user_id, email, otp_code, expires_at, is_verified)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	err := r.conn.QueryRow(ctx, query,
		verification.UserID,
		verification.Email,
		verification.OTPCode,
		verification.ExpiresAt,
		false,
	).Scan(&verification.ID, &verification.CreatedAt)

	return err
}

// GetByEmail retrieves the latest verification record by email
func (r *EmailVerificationRepository) GetByEmail(ctx context.Context, email string) (*models.EmailVerification, error) {
	query := `
		SELECT id, user_id, email, otp_code, expires_at, is_verified, created_at
		FROM email_verifications
		WHERE email = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	verification := &models.EmailVerification{}
	err := r.conn.QueryRow(ctx, query, email).Scan(
		&verification.ID,
		&verification.UserID,
		&verification.Email,
		&verification.OTPCode,
		&verification.ExpiresAt,
		&verification.IsVerified,
		&verification.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return verification, err
}

// MarkAsVerified marks an email verification as verified
func (r *EmailVerificationRepository) MarkAsVerified(ctx context.Context, id int) error {
	query := `UPDATE email_verifications SET is_verified = true WHERE id = $1`
	_, err := r.conn.Exec(ctx, query, id)
	return err
}

// UpdateUserVerification updates user's is_verified status
func (r *EmailVerificationRepository) UpdateUserVerification(ctx context.Context, userID int) error {
	query := `UPDATE users SET is_verified = true WHERE id = $1`
	_, err := r.conn.Exec(ctx, query, userID)
	return err
}

// DeleteExpired deletes expired OTP records (cleanup)
func (r *EmailVerificationRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM email_verifications WHERE expires_at < NOW() AND is_verified = false`
	_, err := r.conn.Exec(ctx, query)
	return err
}
