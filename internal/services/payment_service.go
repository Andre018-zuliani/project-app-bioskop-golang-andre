package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/andre/project-app-bioskop-golang/internal/models"
)

// PaymentService handles payment-related business logic
type PaymentService struct {
	paymentRepo PaymentRepository
	bookingRepo BookingRepository
}

// NewPaymentService creates a new PaymentService
func NewPaymentService(paymentRepo PaymentRepository, bookingRepo BookingRepository) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		bookingRepo: bookingRepo,
	}
}

// ProcessPayment processes a payment for a booking
func (s *PaymentService) ProcessPayment(ctx context.Context, userID int, req *models.PaymentRequest) (*models.PaymentResponse, error) {
	// Get booking
	booking, err := s.bookingRepo.GetBookingByID(ctx, req.BookingID)
	if err != nil {
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}
	if booking == nil {
		return nil, errors.New("booking not found")
	}

	// Verify user owns the booking
	if booking.UserID != userID {
		return nil, errors.New("unauthorized to pay for this booking")
	}

	// Verify amount
	if req.Amount != booking.TotalPrice {
		return nil, fmt.Errorf("amount mismatch: expected %.2f, got %.2f", booking.TotalPrice, req.Amount)
	}

	// Check payment method exists
	method, err := s.paymentRepo.GetPaymentMethodByName(ctx, req.PaymentMethod)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment method: %w", err)
	}
	if method == nil {
		return nil, errors.New("invalid payment method")
	}

	// Create payment
	payment := &models.Payment{
		BookingID:     req.BookingID,
		UserID:        userID,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		Status:        "success", // In real app, this would depend on payment gateway response
		TransactionID: fmt.Sprintf("TXN-%d-%d", booking.ID, booking.UserID),
	}

	err = s.paymentRepo.CreatePayment(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Update booking payment status
	err = s.bookingRepo.UpdateBookingPaymentStatus(ctx, req.BookingID, "paid")
	if err != nil {
		return nil, fmt.Errorf("failed to update booking payment status: %w", err)
	}

	// Update booking status to confirmed
	err = s.bookingRepo.UpdateBookingStatus(ctx, req.BookingID, "confirmed")
	if err != nil {
		return nil, fmt.Errorf("failed to update booking status: %w", err)
	}

	response := &models.PaymentResponse{
		ID:            payment.ID,
		BookingID:     payment.BookingID,
		Amount:        payment.Amount,
		PaymentMethod: payment.PaymentMethod,
		Status:        payment.Status,
		TransactionID: payment.TransactionID,
		CreatedAt:     payment.CreatedAt,
	}

	return response, nil
}

// GetPaymentMethods retrieves all available payment methods
func (s *PaymentService) GetPaymentMethods(ctx context.Context) ([]*models.PaymentMethod, error) {
	methods, err := s.paymentRepo.GetPaymentMethods(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment methods: %w", err)
	}
	return methods, nil
}

// GetPaymentByID retrieves a payment by ID
func (s *PaymentService) GetPaymentByID(ctx context.Context, id int) (*models.Payment, error) {
	payment, err := s.paymentRepo.GetPaymentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	return payment, nil
}
