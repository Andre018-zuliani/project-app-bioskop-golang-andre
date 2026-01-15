package models

import "time"

// Booking represents a seat booking
type Booking struct {
	ID            int       `db:"id" json:"id"`
	UserID        int       `db:"user_id" json:"user_id"`
	CinemaID      int       `db:"cinema_id" json:"cinema_id"`
	SeatID        int       `db:"seat_id" json:"seat_id"`
	ShowDate      time.Time `db:"show_date" json:"show_date"`
	ShowTime      string    `db:"show_time" json:"show_time"`
	BookingDate   time.Time `db:"booking_date" json:"booking_date"`
	Status        string    `db:"status" json:"status"` // pending, confirmed, cancelled
	TotalPrice    float64   `db:"total_price" json:"total_price"`
	PaymentMethod string    `db:"payment_method" json:"payment_method"`
	PaymentStatus string    `db:"payment_status" json:"payment_status"` // pending, paid, failed
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
	Cinema        *Cinema   `json:"cinema,omitempty"`
	Seat          *Seat     `json:"seat,omitempty"`
}

// BookingRequest represents the request body for creating a booking
type BookingRequest struct {
	CinemaID      int    `json:"cinema_id" validate:"required"`
	SeatID        int    `json:"seat_id" validate:"required"`
	Date          string `json:"date" validate:"required"`
	Time          string `json:"time" validate:"required"`
	PaymentMethod string `json:"payment_method" validate:"required"`
}

// BookingResponse represents a booking response
type BookingResponse struct {
	ID            int       `json:"id"`
	CinemaID      int       `json:"cinema_id"`
	SeatID        int       `json:"seat_id"`
	ShowDate      time.Time `json:"show_date"`
	ShowTime      string    `json:"show_time"`
	TotalPrice    float64   `json:"total_price"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	PaymentStatus string    `json:"payment_status"`
	CreatedAt     time.Time `json:"created_at"`
}

// UserBookingHistory represents booking history for a user
type UserBookingHistory struct {
	ID            int       `json:"id"`
	CinemaName    string    `json:"cinema_name"`
	SeatNumber    string    `json:"seat_number"`
	ShowDate      time.Time `json:"show_date"`
	ShowTime      string    `json:"show_time"`
	TotalPrice    float64   `json:"total_price"`
	Status        string    `json:"status"`
	PaymentStatus string    `json:"payment_status"`
	BookingDate   time.Time `json:"booking_date"`
}
