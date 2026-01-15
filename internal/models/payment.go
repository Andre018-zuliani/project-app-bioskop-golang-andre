package models

import "time"

// PaymentMethod represents a payment method
type PaymentMethod struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Type      string    `db:"type" json:"type"` // credit_card, debit_card, e_wallet, transfer
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Payment represents a payment transaction
type Payment struct {
	ID            int       `db:"id" json:"id"`
	BookingID     int       `db:"booking_id" json:"booking_id"`
	UserID        int       `db:"user_id" json:"user_id"`
	Amount        float64   `db:"amount" json:"amount"`
	PaymentMethod string    `db:"payment_method" json:"payment_method"`
	Status        string    `db:"status" json:"status"` // pending, success, failed
	TransactionID string    `db:"transaction_id" json:"transaction_id"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

// PaymentRequest represents the request for payment processing
type PaymentRequest struct {
	BookingID     int     `json:"booking_id" validate:"required"`
	PaymentMethod string  `json:"payment_method" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
}

// PaymentResponse represents a payment response
type PaymentResponse struct {
	ID            int       `json:"id"`
	BookingID     int       `json:"booking_id"`
	Amount        float64   `json:"amount"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	TransactionID string    `json:"transaction_id"`
	CreatedAt     time.Time `json:"created_at"`
}
