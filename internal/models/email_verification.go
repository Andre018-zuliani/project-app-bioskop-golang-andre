package models

import "time"

// EmailVerification represents an email verification record with OTP
type EmailVerification struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Email      string    `json:"email"`
	OTPCode    string    `json:"-"` // Hidden from JSON responses
	ExpiresAt  time.Time `json:"expires_at"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
}

// SendOTPRequest represents request to send OTP email
type SendOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// VerifyOTPRequest represents request to verify OTP
type VerifyOTPRequest struct {
	Email   string `json:"email" validate:"required,email"`
	OTPCode string `json:"otp_code" validate:"required,len=6"`
}

// OTPResponse represents OTP verification response
type OTPResponse struct {
	Message    string `json:"message"`
	IsVerified bool   `json:"is_verified"`
}
