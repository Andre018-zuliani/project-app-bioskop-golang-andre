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

// MockBookingRepository is a mock implementation of BookingRepository
type MockBookingRepository struct {
	mock.Mock
}

func (m *MockBookingRepository) CreateBooking(ctx context.Context, booking *models.Booking) error {
	args := m.Called(ctx, booking)
	return args.Error(0)
}

func (m *MockBookingRepository) GetBookingByID(ctx context.Context, id int) (*models.Booking, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Booking), args.Error(1)
}

func (m *MockBookingRepository) GetUserBookings(ctx context.Context, userID, page, limit int) ([]*models.Booking, int, error) {
	args := m.Called(ctx, userID, page, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*models.Booking), args.Int(1), args.Error(2)
}

func (m *MockBookingRepository) UpdateBookingStatus(ctx context.Context, id int, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockBookingRepository) UpdateBookingPaymentStatus(ctx context.Context, id int, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockBookingRepository) CheckSeatBooked(ctx context.Context, seatID int, date time.Time, showTime string) (bool, error) {
	args := m.Called(ctx, seatID, date, showTime)
	return args.Bool(0), args.Error(1)
}

func (m *MockBookingRepository) GetBookingWithDetails(ctx context.Context, id int) (*models.Booking, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Booking), args.Error(1)
}

// MockSeatRepository is a mock implementation of SeatRepository
type MockSeatRepository struct {
	mock.Mock
}

func (m *MockSeatRepository) GetSeatsByCinema(ctx context.Context, cinemaID int) ([]*models.Seat, error) {
	args := m.Called(ctx, cinemaID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Seat), args.Error(1)
}

func (m *MockSeatRepository) GetSeatByID(ctx context.Context, id int) (*models.Seat, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Seat), args.Error(1)
}

func (m *MockSeatRepository) CreateSeat(ctx context.Context, seat *models.Seat) error {
	args := m.Called(ctx, seat)
	return args.Error(0)
}

func (m *MockSeatRepository) GetSeatAvailability(ctx context.Context, cinemaID int, date time.Time, timeStr string) ([]*models.SeatAvailability, error) {
	args := m.Called(ctx, cinemaID, date, timeStr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.SeatAvailability), args.Error(1)
}

func (m *MockSeatRepository) UpdateSeatAvailability(ctx context.Context, seatID int, date time.Time, timeStr string, isAvailable bool) error {
	args := m.Called(ctx, seatID, date, timeStr, isAvailable)
	return args.Error(0)
}

// MockCinemaRepository is a mock implementation of CinemaRepository
type MockCinemaRepository struct {
	mock.Mock
}

func (m *MockCinemaRepository) GetAllCinemas(ctx context.Context, page, limit int, filters *models.CinemaFilters) ([]*models.Cinema, int, error) {
	args := m.Called(ctx, page, limit, filters)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*models.Cinema), args.Int(1), args.Error(2)
}

func (m *MockCinemaRepository) GetCinemaByID(ctx context.Context, id int) (*models.Cinema, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Cinema), args.Error(1)
}

func (m *MockCinemaRepository) CreateCinema(ctx context.Context, cinema *models.Cinema) error {
	args := m.Called(ctx, cinema)
	return args.Error(0)
}

// TestCreateBooking_Success tests successful booking creation
func TestCreateBooking_Success(t *testing.T) {
	// Arrange
	mockBookingRepo := new(MockBookingRepository)
	mockSeatRepo := new(MockSeatRepository)
	mockCinemaRepo := new(MockCinemaRepository)
	service := NewBookingService(mockBookingRepo, mockSeatRepo, mockCinemaRepo)

	userID := 1
	req := &models.BookingRequest{
		CinemaID:      1,
		SeatID:        1,
		Date:          "2026-01-15",
		Time:          "19:00",
		PaymentMethod: "credit_card",
	}

	showDate, _ := time.Parse("2006-01-02", req.Date)

	seat := &models.Seat{ID: 1, CinemaID: 1, SeatNumber: "A1", Price: 50000}
	cinema := &models.Cinema{ID: 1, Name: "Cinema XXI", City: "Jakarta"}

	mockSeatRepo.On("GetSeatByID", mock.Anything, 1).Return(seat, nil)
	mockCinemaRepo.On("GetCinemaByID", mock.Anything, 1).Return(cinema, nil)
	mockBookingRepo.On("CheckSeatBooked", mock.Anything, 1, showDate, "19:00").Return(false, nil)
	mockBookingRepo.On("CreateBooking", mock.Anything, mock.AnythingOfType("*models.Booking")).Run(func(args mock.Arguments) {
		b := args.Get(1).(*models.Booking)
		b.ID = 1
		b.CreatedAt = time.Now()
	}).Return(nil)
	mockSeatRepo.On("UpdateSeatAvailability", mock.Anything, 1, showDate, "19:00", false).Return(nil)

	// Act
	response, err := service.CreateBooking(context.Background(), userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 1, response.ID)
	mockBookingRepo.AssertExpectations(t)
	mockSeatRepo.AssertExpectations(t)
	mockCinemaRepo.AssertExpectations(t)
}

func TestCreateBooking_SeatNotFound(t *testing.T) {
	// Arrange
	mockBookingRepo := new(MockBookingRepository)
	mockSeatRepo := new(MockSeatRepository)
	mockCinemaRepo := new(MockCinemaRepository)
	service := NewBookingService(mockBookingRepo, mockSeatRepo, mockCinemaRepo)

	userID := 1
	req := &models.BookingRequest{
		CinemaID:      1,
		SeatID:        999,
		Date:          "2026-01-15",
		Time:          "19:00",
		PaymentMethod: "credit_card",
	}

	mockSeatRepo.On("GetSeatByID", mock.Anything, 999).Return(nil, errors.New("seat not found"))

	// Act
	response, err := service.CreateBooking(context.Background(), userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "seat")
	mockSeatRepo.AssertExpectations(t)
}

func TestCreateBooking_SeatAlreadyBooked(t *testing.T) {
	// Arrange
	mockBookingRepo := new(MockBookingRepository)
	mockSeatRepo := new(MockSeatRepository)
	mockCinemaRepo := new(MockCinemaRepository)
	service := NewBookingService(mockBookingRepo, mockSeatRepo, mockCinemaRepo)

	userID := 1
	req := &models.BookingRequest{
		CinemaID:      1,
		SeatID:        1,
		Date:          "2026-01-15",
		Time:          "19:00",
		PaymentMethod: "credit_card",
	}

	seat := &models.Seat{ID: 1, CinemaID: 1, SeatNumber: "A1", Price: 50000}
	cinema := &models.Cinema{ID: 1, Name: "Cinema XXI", City: "Jakarta"}

	showDate, _ := time.Parse("2006-01-02", req.Date)

	mockSeatRepo.On("GetSeatByID", mock.Anything, 1).Return(seat, nil)
	mockCinemaRepo.On("GetCinemaByID", mock.Anything, 1).Return(cinema, nil)
	mockBookingRepo.On("CheckSeatBooked", mock.Anything, 1, showDate, "19:00").Return(true, nil)

	// Act
	response, err := service.CreateBooking(context.Background(), userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "already booked")
	mockBookingRepo.AssertExpectations(t)
	mockSeatRepo.AssertExpectations(t)
	mockCinemaRepo.AssertExpectations(t)
}

// TestGetUserBookings_Success tests retrieving user bookings
func TestGetUserBookings_Success(t *testing.T) {
	// Arrange
	mockBookingRepo := new(MockBookingRepository)
	mockSeatRepo := new(MockSeatRepository)
	mockCinemaRepo := new(MockCinemaRepository)
	service := NewBookingService(mockBookingRepo, mockSeatRepo, mockCinemaRepo)

	userID := 1
	page := 1
	limit := 10

	showDate := time.Now()
	baseBookings := []*models.Booking{{ID: 1, UserID: userID, ShowDate: showDate, ShowTime: "19:00"}}
	detailed := &models.Booking{ID: 1, UserID: userID, ShowDate: showDate, ShowTime: "19:00", Status: "confirmed"}

	mockBookingRepo.On("GetUserBookings", mock.Anything, userID, page, limit).Return(baseBookings, 1, nil)
	mockBookingRepo.On("GetBookingWithDetails", mock.Anything, 1).Return(detailed, nil)

	// Act
	response, err := service.GetUserBookings(context.Background(), userID, page, limit)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 1, len(response.Data.([]*models.Booking)))
	assert.Equal(t, 1, response.Total)
	mockBookingRepo.AssertExpectations(t)
}

func TestGetUserBookings_EmptyResult(t *testing.T) {
	// Arrange
	mockBookingRepo := new(MockBookingRepository)
	mockSeatRepo := new(MockSeatRepository)
	mockCinemaRepo := new(MockCinemaRepository)
	service := NewBookingService(mockBookingRepo, mockSeatRepo, mockCinemaRepo)

	userID := 1
	page := 1
	limit := 10

	mockBookingRepo.On("GetUserBookings", mock.Anything, userID, page, limit).Return([]*models.Booking{}, 0, nil)

	// Act
	response, err := service.GetUserBookings(context.Background(), userID, page, limit)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 0, len(response.Data.([]*models.Booking)))
	assert.Equal(t, 0, response.Total)
	mockBookingRepo.AssertExpectations(t)
}

// TestUpdateBookingStatus tests booking status update
func TestUpdateBookingStatus_Success(t *testing.T) {
	// Arrange
	mockBookingRepo := new(MockBookingRepository)
	mockSeatRepo := new(MockSeatRepository)
	mockCinemaRepo := new(MockCinemaRepository)
	service := NewBookingService(mockBookingRepo, mockSeatRepo, mockCinemaRepo)

	bookingID := 1
	newStatus := "confirmed"

	mockBookingRepo.On("UpdateBookingStatus", mock.Anything, bookingID, newStatus).Return(nil)

	// Act
	err := service.UpdateBookingStatus(context.Background(), bookingID, newStatus)

	// Assert
	assert.NoError(t, err)
	mockBookingRepo.AssertExpectations(t)
}

func TestCreateBooking_InvalidDate(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockSeatRepo := new(MockSeatRepository)
	mockCinemaRepo := new(MockCinemaRepository)
	service := NewBookingService(mockBookingRepo, mockSeatRepo, mockCinemaRepo)

	userID := 1
	req := &models.BookingRequest{Date: "15-01-2026", Time: "19:00"}

	response, err := service.CreateBooking(context.Background(), userID, req)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "invalid date")
}

func TestCreateBooking_CinemaNotFound(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockSeatRepo := new(MockSeatRepository)
	mockCinemaRepo := new(MockCinemaRepository)
	service := NewBookingService(mockBookingRepo, mockSeatRepo, mockCinemaRepo)

	req := &models.BookingRequest{CinemaID: 10, SeatID: 1, Date: "2026-01-15", Time: "19:00"}
	showDate, _ := time.Parse("2006-01-02", req.Date)

	mockSeatRepo.On("GetSeatByID", mock.Anything, 1).Return(&models.Seat{ID: 1, CinemaID: 10, Price: 50000}, nil)
	mockCinemaRepo.On("GetCinemaByID", mock.Anything, 10).Return(nil, nil)

	resp, err := service.CreateBooking(context.Background(), 1, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "cinema not found")
	mockSeatRepo.AssertExpectations(t)
	mockCinemaRepo.AssertExpectations(t)

	_ = showDate // keep showDate used to mirror other tests
}

func TestCreateBooking_SeatWrongCinema(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockSeatRepo := new(MockSeatRepository)
	mockCinemaRepo := new(MockCinemaRepository)
	service := NewBookingService(mockBookingRepo, mockSeatRepo, mockCinemaRepo)

	req := &models.BookingRequest{CinemaID: 2, SeatID: 1, Date: "2026-01-15", Time: "19:00"}

	mockSeatRepo.On("GetSeatByID", mock.Anything, 1).Return(&models.Seat{ID: 1, CinemaID: 1, Price: 50000}, nil)
	mockCinemaRepo.On("GetCinemaByID", mock.Anything, 2).Return(&models.Cinema{ID: 2}, nil)

	resp, err := service.CreateBooking(context.Background(), 1, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "does not belong")
	mockSeatRepo.AssertExpectations(t)
	mockCinemaRepo.AssertExpectations(t)
}

func TestCreateBooking_CheckSeatBookedError(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockSeatRepo := new(MockSeatRepository)
	mockCinemaRepo := new(MockCinemaRepository)
	service := NewBookingService(mockBookingRepo, mockSeatRepo, mockCinemaRepo)

	req := &models.BookingRequest{CinemaID: 1, SeatID: 1, Date: "2026-01-15", Time: "19:00"}
	showDate, _ := time.Parse("2006-01-02", req.Date)

	mockSeatRepo.On("GetSeatByID", mock.Anything, 1).Return(&models.Seat{ID: 1, CinemaID: 1, Price: 50000}, nil)
	mockCinemaRepo.On("GetCinemaByID", mock.Anything, 1).Return(&models.Cinema{ID: 1}, nil)
	mockBookingRepo.On("CheckSeatBooked", mock.Anything, 1, showDate, "19:00").Return(false, errors.New("db error"))

	resp, err := service.CreateBooking(context.Background(), 1, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to check seat booking")
	mockBookingRepo.AssertExpectations(t)
}

func TestCreateBooking_CreateBookingError(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockSeatRepo := new(MockSeatRepository)
	mockCinemaRepo := new(MockCinemaRepository)
	service := NewBookingService(mockBookingRepo, mockSeatRepo, mockCinemaRepo)

	req := &models.BookingRequest{CinemaID: 1, SeatID: 1, Date: "2026-01-15", Time: "19:00", PaymentMethod: "cash"}
	showDate, _ := time.Parse("2006-01-02", req.Date)

	mockSeatRepo.On("GetSeatByID", mock.Anything, 1).Return(&models.Seat{ID: 1, CinemaID: 1, Price: 50000}, nil)
	mockCinemaRepo.On("GetCinemaByID", mock.Anything, 1).Return(&models.Cinema{ID: 1}, nil)
	mockBookingRepo.On("CheckSeatBooked", mock.Anything, 1, showDate, "19:00").Return(false, nil)
	mockBookingRepo.On("CreateBooking", mock.Anything, mock.AnythingOfType("*models.Booking")).Return(errors.New("insert fail"))

	resp, err := service.CreateBooking(context.Background(), 1, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "create booking")
	mockBookingRepo.AssertExpectations(t)
}

func TestCreateBooking_UpdateAvailabilityError(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockSeatRepo := new(MockSeatRepository)
	mockCinemaRepo := new(MockCinemaRepository)
	service := NewBookingService(mockBookingRepo, mockSeatRepo, mockCinemaRepo)

	req := &models.BookingRequest{CinemaID: 1, SeatID: 1, Date: "2026-01-15", Time: "19:00", PaymentMethod: "cash"}
	showDate, _ := time.Parse("2006-01-02", req.Date)

	mockSeatRepo.On("GetSeatByID", mock.Anything, 1).Return(&models.Seat{ID: 1, CinemaID: 1, Price: 50000}, nil)
	mockCinemaRepo.On("GetCinemaByID", mock.Anything, 1).Return(&models.Cinema{ID: 1}, nil)
	mockBookingRepo.On("CheckSeatBooked", mock.Anything, 1, showDate, "19:00").Return(false, nil)
	mockBookingRepo.On("CreateBooking", mock.Anything, mock.AnythingOfType("*models.Booking")).Return(nil)
	mockSeatRepo.On("UpdateSeatAvailability", mock.Anything, 1, showDate, "19:00", false).Return(errors.New("update fail"))

	resp, err := service.CreateBooking(context.Background(), 1, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "update seat availability")
	mockSeatRepo.AssertExpectations(t)
}

func TestGetUserBookings_RepoError(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockSeatRepo := new(MockSeatRepository)
	mockCinemaRepo := new(MockCinemaRepository)
	service := NewBookingService(mockBookingRepo, mockSeatRepo, mockCinemaRepo)

	mockBookingRepo.On("GetUserBookings", mock.Anything, 1, 1, 10).Return(nil, 0, errors.New("query fail"))

	resp, err := service.GetUserBookings(context.Background(), 1, 1, 10)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockBookingRepo.AssertExpectations(t)
}

func TestGetUserBookings_DetailError(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockSeatRepo := new(MockSeatRepository)
	mockCinemaRepo := new(MockCinemaRepository)
	service := NewBookingService(mockBookingRepo, mockSeatRepo, mockCinemaRepo)

	bookings := []*models.Booking{{ID: 1}}
	mockBookingRepo.On("GetUserBookings", mock.Anything, 1, 1, 10).Return(bookings, 1, nil)
	mockBookingRepo.On("GetBookingWithDetails", mock.Anything, 1).Return(nil, errors.New("detail fail"))

	resp, err := service.GetUserBookings(context.Background(), 1, 1, 10)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockBookingRepo.AssertExpectations(t)
}

func TestGetBookingByID_Success(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockSeatRepo := new(MockSeatRepository)
	mockCinemaRepo := new(MockCinemaRepository)
	service := NewBookingService(mockBookingRepo, mockSeatRepo, mockCinemaRepo)

	booking := &models.Booking{ID: 7}
	mockBookingRepo.On("GetBookingWithDetails", mock.Anything, 7).Return(booking, nil)

	result, err := service.GetBookingByID(context.Background(), 7)

	assert.NoError(t, err)
	assert.Equal(t, 7, result.ID)
	mockBookingRepo.AssertExpectations(t)
}

func TestGetBookingByID_Error(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockSeatRepo := new(MockSeatRepository)
	mockCinemaRepo := new(MockCinemaRepository)
	service := NewBookingService(mockBookingRepo, mockSeatRepo, mockCinemaRepo)

	mockBookingRepo.On("GetBookingWithDetails", mock.Anything, 7).Return(nil, errors.New("db fail"))

	result, err := service.GetBookingByID(context.Background(), 7)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get booking")
	mockBookingRepo.AssertExpectations(t)
}
