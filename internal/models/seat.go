package models

import "time"

// Seat represents a cinema seat
type Seat struct {
	ID         int       `db:"id" json:"id"`
	CinemaID   int       `db:"cinema_id" json:"cinema_id"`
	SeatNumber string    `db:"seat_number" json:"seat_number"`
	RowNumber  int       `db:"row_number" json:"row_number"`
	SeatType   string    `db:"seat_type" json:"seat_type"` // standard, premium, vip
	Price      float64   `db:"price" json:"price"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// SeatAvailability represents seat availability for a specific date and time
type SeatAvailability struct {
	ID          int       `db:"id" json:"id"`
	CinemaID    int       `db:"cinema_id" json:"cinema_id"`
	SeatID      int       `db:"seat_id" json:"seat_id"`
	ShowDate    time.Time `db:"show_date" json:"show_date"`
	ShowTime    string    `db:"show_time" json:"show_time"`
	IsAvailable bool      `db:"is_available" json:"is_available"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	Seat        *Seat     `json:"seat,omitempty"`
}

// SeatAvailabilityRequest represents the request for checking seat availability
type SeatAvailabilityRequest struct {
	CinemaID int    `json:"cinema_id" validate:"required"`
	Date     string `json:"date" validate:"required"`
	Time     string `json:"time" validate:"required"`
}

// SeatAvailabilityResponse represents the response for seat availability
type SeatAvailabilityResponse struct {
	CinemaID         int                 `json:"cinema_id"`
	Date             string              `json:"date"`
	Time             string              `json:"time"`
	AvailableSeats   []*SeatAvailability `json:"available_seats"`
	UnavailableSeats []*SeatAvailability `json:"unavailable_seats"`
	TotalAvailable   int                 `json:"total_available"`
	TotalUnavailable int                 `json:"total_unavailable"`
}
