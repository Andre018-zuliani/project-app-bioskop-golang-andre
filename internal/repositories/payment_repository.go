package repositories

import (
	"context"
	"fmt"

	"github.com/andre/project-app-bioskop-golang/internal/models"
	"github.com/jackc/pgx/v5"
)

// PaymentRepository handles payment-related database operations
type PaymentRepository struct {
	db Database
}

// NewPaymentRepository creates a new PaymentRepository
func NewPaymentRepository(db Database) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// CreatePayment creates a new payment
func (r *PaymentRepository) CreatePayment(ctx context.Context, payment *models.Payment) error {
	query := `INSERT INTO payments (booking_id, user_id, amount, payment_method, status, transaction_id) 
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, payment.BookingID, payment.UserID, payment.Amount, payment.PaymentMethod,
		payment.Status, payment.TransactionID).
		Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}
	return nil
}

// GetPaymentByID retrieves a payment by ID
func (r *PaymentRepository) GetPaymentByID(ctx context.Context, id int) (*models.Payment, error) {
	payment := &models.Payment{}
	query := `SELECT id, booking_id, user_id, amount, payment_method, status, transaction_id, created_at, updated_at 
	FROM payments WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).
		Scan(&payment.ID, &payment.BookingID, &payment.UserID, &payment.Amount, &payment.PaymentMethod,
			&payment.Status, &payment.TransactionID, &payment.CreatedAt, &payment.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return payment, nil
}

// GetPaymentByBookingID retrieves a payment by booking ID
func (r *PaymentRepository) GetPaymentByBookingID(ctx context.Context, bookingID int) (*models.Payment, error) {
	payment := &models.Payment{}
	query := `SELECT id, booking_id, user_id, amount, payment_method, status, transaction_id, created_at, updated_at 
	FROM payments WHERE booking_id = $1`

	err := r.db.QueryRow(ctx, query, bookingID).
		Scan(&payment.ID, &payment.BookingID, &payment.UserID, &payment.Amount, &payment.PaymentMethod,
			&payment.Status, &payment.TransactionID, &payment.CreatedAt, &payment.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get payment by booking id: %w", err)
	}

	return payment, nil
}

// UpdatePaymentStatus updates the status of a payment
func (r *PaymentRepository) UpdatePaymentStatus(ctx context.Context, id int, status string) error {
	query := `UPDATE payments SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}
	return nil
}

// GetPaymentMethods retrieves all active payment methods
func (r *PaymentRepository) GetPaymentMethods(ctx context.Context) ([]*models.PaymentMethod, error) {
	query := `SELECT id, name, type, is_active, created_at, updated_at FROM payment_methods WHERE is_active = TRUE ORDER BY name`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment methods: %w", err)
	}
	defer rows.Close()

	methods := []*models.PaymentMethod{}
	for rows.Next() {
		method := &models.PaymentMethod{}
		err := rows.Scan(&method.ID, &method.Name, &method.Type, &method.IsActive, &method.CreatedAt, &method.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment method: %w", err)
		}
		methods = append(methods, method)
	}

	return methods, nil
}

// GetPaymentMethodByName retrieves a payment method by name
func (r *PaymentRepository) GetPaymentMethodByName(ctx context.Context, name string) (*models.PaymentMethod, error) {
	method := &models.PaymentMethod{}
	query := `SELECT id, name, type, is_active, created_at, updated_at FROM payment_methods WHERE name = $1 AND is_active = TRUE`

	err := r.db.QueryRow(ctx, query, name).
		Scan(&method.ID, &method.Name, &method.Type, &method.IsActive, &method.CreatedAt, &method.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get payment method: %w", err)
	}

	return method, nil
}

// GetUserPayments retrieves all payments for a user
func (r *PaymentRepository) GetUserPayments(ctx context.Context, userID int) ([]*models.Payment, error) {
	query := `SELECT id, booking_id, user_id, amount, payment_method, status, transaction_id, created_at, updated_at 
	FROM payments WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user payments: %w", err)
	}
	defer rows.Close()

	payments := []*models.Payment{}
	for rows.Next() {
		payment := &models.Payment{}
		err := rows.Scan(&payment.ID, &payment.BookingID, &payment.UserID, &payment.Amount, &payment.PaymentMethod,
			&payment.Status, &payment.TransactionID, &payment.CreatedAt, &payment.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}
		payments = append(payments, payment)
	}

	return payments, nil
}
