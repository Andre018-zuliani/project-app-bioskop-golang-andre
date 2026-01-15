package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/andre/project-app-bioskop-golang/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) CreatePayment(ctx context.Context, payment *models.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) GetPaymentByID(ctx context.Context, id int) (*models.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Payment), args.Error(1)
}

func (m *MockPaymentRepository) GetPaymentByBookingID(ctx context.Context, bookingID int) (*models.Payment, error) {
	args := m.Called(ctx, bookingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Payment), args.Error(1)
}

func (m *MockPaymentRepository) UpdatePaymentStatus(ctx context.Context, id int, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockPaymentRepository) GetPaymentMethods(ctx context.Context) ([]*models.PaymentMethod, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.PaymentMethod), args.Error(1)
}

func (m *MockPaymentRepository) GetPaymentMethodByName(ctx context.Context, name string) (*models.PaymentMethod, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PaymentMethod), args.Error(1)
}

type MockBookingRepoForPayment struct {
	mock.Mock
}

func (m *MockBookingRepoForPayment) CreateBooking(ctx context.Context, booking *models.Booking) error {
	return errors.New("not implemented")
}

func (m *MockBookingRepoForPayment) GetBookingByID(ctx context.Context, id int) (*models.Booking, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Booking), args.Error(1)
}

func (m *MockBookingRepoForPayment) GetBookingWithDetails(ctx context.Context, id int) (*models.Booking, error) {
	return nil, errors.New("not implemented")
}

func (m *MockBookingRepoForPayment) GetUserBookings(ctx context.Context, userID, page, limit int) ([]*models.Booking, int, error) {
	return nil, 0, errors.New("not implemented")
}

func (m *MockBookingRepoForPayment) UpdateBookingStatus(ctx context.Context, id int, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockBookingRepoForPayment) UpdateBookingPaymentStatus(ctx context.Context, id int, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockBookingRepoForPayment) CheckSeatBooked(ctx context.Context, seatID int, showDate time.Time, showTime string) (bool, error) {
	return false, errors.New("not implemented")
}

func TestProcessPayment_Success(t *testing.T) {
	paymentRepo := new(MockPaymentRepository)
	bookingRepo := new(MockBookingRepoForPayment)
	service := NewPaymentService(paymentRepo, bookingRepo)

	booking := &models.Booking{ID: 1, UserID: 1, TotalPrice: 100000}
	method := &models.PaymentMethod{Name: "Card"}
	req := &models.PaymentRequest{BookingID: 1, Amount: 100000, PaymentMethod: "Card"}

	bookingRepo.On("GetBookingByID", mock.Anything, 1).Return(booking, nil)
	paymentRepo.On("GetPaymentMethodByName", mock.Anything, "Card").Return(method, nil)
	paymentRepo.On("CreatePayment", mock.Anything, mock.AnythingOfType("*models.Payment")).Return(nil)
	bookingRepo.On("UpdateBookingPaymentStatus", mock.Anything, 1, "paid").Return(nil)
	bookingRepo.On("UpdateBookingStatus", mock.Anything, 1, "confirmed").Return(nil)

	resp, err := service.ProcessPayment(context.Background(), 1, req)

	assert.NoError(t, err)
	assert.Equal(t, req.Amount, resp.Amount)
	assert.Equal(t, req.PaymentMethod, resp.PaymentMethod)
	bookingRepo.AssertExpectations(t)
	paymentRepo.AssertExpectations(t)
}

func TestProcessPayment_BookingNotFound(t *testing.T) {
	paymentRepo := new(MockPaymentRepository)
	bookingRepo := new(MockBookingRepoForPayment)
	service := NewPaymentService(paymentRepo, bookingRepo)

	req := &models.PaymentRequest{BookingID: 99, Amount: 50000, PaymentMethod: "Card"}
	bookingRepo.On("GetBookingByID", mock.Anything, 99).Return(nil, nil)

	resp, err := service.ProcessPayment(context.Background(), 1, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestProcessPayment_Unauthorized(t *testing.T) {
	paymentRepo := new(MockPaymentRepository)
	bookingRepo := new(MockBookingRepoForPayment)
	service := NewPaymentService(paymentRepo, bookingRepo)

	booking := &models.Booking{ID: 1, UserID: 2, TotalPrice: 100000}
	req := &models.PaymentRequest{BookingID: 1, Amount: 100000, PaymentMethod: "Card"}

	bookingRepo.On("GetBookingByID", mock.Anything, 1).Return(booking, nil)

	resp, err := service.ProcessPayment(context.Background(), 1, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestProcessPayment_AmountMismatch(t *testing.T) {
	paymentRepo := new(MockPaymentRepository)
	bookingRepo := new(MockBookingRepoForPayment)
	service := NewPaymentService(paymentRepo, bookingRepo)

	booking := &models.Booking{ID: 1, UserID: 1, TotalPrice: 100000}
	req := &models.PaymentRequest{BookingID: 1, Amount: 200000, PaymentMethod: "Card"}

	bookingRepo.On("GetBookingByID", mock.Anything, 1).Return(booking, nil)

	resp, err := service.ProcessPayment(context.Background(), 1, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestProcessPayment_InvalidMethod(t *testing.T) {
	paymentRepo := new(MockPaymentRepository)
	bookingRepo := new(MockBookingRepoForPayment)
	service := NewPaymentService(paymentRepo, bookingRepo)

	booking := &models.Booking{ID: 1, UserID: 1, TotalPrice: 100000}
	req := &models.PaymentRequest{BookingID: 1, Amount: 100000, PaymentMethod: "Unknown"}

	bookingRepo.On("GetBookingByID", mock.Anything, 1).Return(booking, nil)
	paymentRepo.On("GetPaymentMethodByName", mock.Anything, "Unknown").Return(nil, nil)

	resp, err := service.ProcessPayment(context.Background(), 1, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestGetPaymentMethods(t *testing.T) {
	paymentRepo := new(MockPaymentRepository)
	bookingRepo := new(MockBookingRepoForPayment)
	service := NewPaymentService(paymentRepo, bookingRepo)

	methods := []*models.PaymentMethod{{ID: 1, Name: "Card"}}
	paymentRepo.On("GetPaymentMethods", mock.Anything).Return(methods, nil)

	result, err := service.GetPaymentMethods(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestGetPaymentByID(t *testing.T) {
	paymentRepo := new(MockPaymentRepository)
	bookingRepo := new(MockBookingRepoForPayment)
	service := NewPaymentService(paymentRepo, bookingRepo)

	payment := &models.Payment{ID: 10}
	paymentRepo.On("GetPaymentByID", mock.Anything, 10).Return(payment, nil)

	result, err := service.GetPaymentByID(context.Background(), 10)

	assert.NoError(t, err)
	assert.Equal(t, payment, result)
}
