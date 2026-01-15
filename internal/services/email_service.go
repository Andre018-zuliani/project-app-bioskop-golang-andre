package services

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/andre/project-app-bioskop-golang/internal/models"
	"github.com/andre/project-app-bioskop-golang/internal/repositories"
	"go.uber.org/zap"
)

// EmailService handles email OTP operations
type EmailService struct {
	emailRepo *repositories.EmailVerificationRepository
	logger    *zap.Logger
	apiURL    string
	apiKey    string
}

// NewEmailService creates a new email service
func NewEmailService(emailRepo *repositories.EmailVerificationRepository, logger *zap.Logger, apiURL, apiKey string) *EmailService {
	return &EmailService{
		emailRepo: emailRepo,
		logger:    logger,
		apiURL:    apiURL,
		apiKey:    apiKey,
	}
}

// GenerateOTP generates a secure 6-digit OTP
func (s *EmailService) GenerateOTP() (string, error) {
	otp := ""
	for i := 0; i < 6; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		otp += num.String()
	}
	return otp, nil
}

// SendOTP generates and sends OTP to user's email
func (s *EmailService) SendOTP(ctx context.Context, userID int, email, username string) error {
	// Generate OTP
	otpCode, err := s.GenerateOTP()
	if err != nil {
		s.logger.Error("Failed to generate OTP", zap.Error(err))
		return errors.New("failed to generate OTP")
	}

	// Save OTP to database (expires in 5 minutes)
	verification := &models.EmailVerification{
		UserID:    userID,
		Email:     email,
		OTPCode:   otpCode,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	err = s.emailRepo.Create(ctx, verification)
	if err != nil {
		s.logger.Error("Failed to save OTP", zap.Error(err))
		return errors.New("failed to save OTP")
	}

	// Send email asynchronously (non-blocking)
	go func() {
		err := s.sendEmailViaAPI(email, username, otpCode)
		if err != nil {
			s.logger.Error("Failed to send email", zap.Error(err), zap.String("email", email))
		} else {
			s.logger.Info("OTP email sent successfully", zap.String("email", email))
		}
	}()

	return nil
}

// sendEmailViaAPI sends email via Lumoshive Email API
func (s *EmailService) sendEmailViaAPI(toEmail, name, otpCode string) error {
	emailBody := fmt.Sprintf(`
Halo %s,

Kode OTP untuk verifikasi email Anda adalah:

%s

Kode ini berlaku selama 5 menit.

Jika Anda tidak merasa melakukan registrasi, abaikan email ini.

Terima kasih,
Cinema Booking System Team
`, name, otpCode)

	reqBody := map[string]string{
		"to":      toEmail,
		"name":    name,
		"subject": "Cinema Booking System - Kode OTP Verifikasi Email",
		"text":    emailBody,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("email API returned status: %d", resp.StatusCode)
	}

	return nil
}

// VerifyOTP verifies the OTP code provided by user
func (s *EmailService) VerifyOTP(ctx context.Context, email, otpCode string) error {
	// Get latest verification record
	verification, err := s.emailRepo.GetByEmail(ctx, email)
	if err != nil {
		s.logger.Error("Failed to get verification", zap.Error(err))
		return errors.New("verification record not found")
	}

	if verification == nil {
		return errors.New("no OTP found for this email")
	}

	// Check if already verified
	if verification.IsVerified {
		return errors.New("email already verified")
	}

	// Check if expired
	if time.Now().After(verification.ExpiresAt) {
		return errors.New("OTP has expired, please request a new one")
	}

	// Verify OTP code
	if verification.OTPCode != otpCode {
		return errors.New("invalid OTP code")
	}

	// Mark as verified
	err = s.emailRepo.MarkAsVerified(ctx, verification.ID)
	if err != nil {
		s.logger.Error("Failed to mark as verified", zap.Error(err))
		return errors.New("failed to verify email")
	}

	// Update user's is_verified status
	err = s.emailRepo.UpdateUserVerification(ctx, verification.UserID)
	if err != nil {
		s.logger.Error("Failed to update user verification", zap.Error(err))
		return errors.New("failed to update user status")
	}

	s.logger.Info("Email verified successfully", zap.String("email", email))
	return nil
}

// ResendOTP resends OTP to user's email
func (s *EmailService) ResendOTP(ctx context.Context, email string) error {
	// Get user's latest verification
	verification, err := s.emailRepo.GetByEmail(ctx, email)
	if err != nil || verification == nil {
		return errors.New("no registration found for this email")
	}

	if verification.IsVerified {
		return errors.New("email already verified")
	}

	// Check rate limiting (allow resend after 1 minute)
	if time.Since(verification.CreatedAt) < 1*time.Minute {
		return errors.New("please wait 1 minute before requesting a new OTP")
	}

	// Generate new OTP
	otpCode, err := s.GenerateOTP()
	if err != nil {
		return errors.New("failed to generate new OTP")
	}

	// Create new verification record
	newVerification := &models.EmailVerification{
		UserID:    verification.UserID,
		Email:     email,
		OTPCode:   otpCode,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	err = s.emailRepo.Create(ctx, newVerification)
	if err != nil {
		s.logger.Error("Failed to save new OTP", zap.Error(err))
		return errors.New("failed to save new OTP")
	}

	// Send email asynchronously
	go func() {
		// Get username from verification (we need to store it or fetch from users table)
		// For now, use email as name
		err := s.sendEmailViaAPI(email, email, otpCode)
		if err != nil {
			s.logger.Error("Failed to resend email", zap.Error(err))
		}
	}()

	s.logger.Info("OTP resent successfully", zap.String("email", email))
	return nil
}
