package models

import "time"

// Cinema represents a cinema venue
type Cinema struct {
	ID         int       `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	Location   string    `db:"location" json:"location"`
	City       string    `db:"city" json:"city"`
	Address    string    `db:"address" json:"address"`
	TotalSeats int       `db:"total_seats" json:"total_seats"`
	ImageURL   string    `db:"image_url" json:"image_url"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// CinemaFilters represents filters for cinema listing
type CinemaFilters struct {
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	City  string `json:"city"`
	Name  string `json:"name"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int         `json:"total"`
	TotalPages int         `json:"total_pages"`
}
