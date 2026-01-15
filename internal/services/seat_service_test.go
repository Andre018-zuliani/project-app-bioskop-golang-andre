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

type MockSeatAvailabilityRepository struct {
	mock.Mock
}

func (m *MockSeatAvailabilityRepository) GetSeatAvailability(ctx context.Context, cinemaID int, date time.Time, timeStr string) ([]*models.SeatAvailability, error) {
	args := m.Called(ctx, cinemaID, date, timeStr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.SeatAvailability), args.Error(1)
}

func (m *MockSeatAvailabilityRepository) GetSeatByID(ctx context.Context, id int) (*models.Seat, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Seat), args.Error(1)
}

func (m *MockSeatAvailabilityRepository) UpdateSeatAvailability(ctx context.Context, seatID int, date time.Time, timeStr string, isAvailable bool) error {
	args := m.Called(ctx, seatID, date, timeStr, isAvailable)
	return args.Error(0)
}

func TestGetSeatAvailability_Success(t *testing.T) {
	repo := new(MockSeatAvailabilityRepository)
	service := NewSeatService(repo)

	dateStr := "2026-01-15"
	showDate, _ := time.Parse("2006-01-02", dateStr)

	availabilities := []*models.SeatAvailability{
		{ID: 1, CinemaID: 1, SeatID: 1, ShowDate: showDate, ShowTime: "19:00", IsAvailable: true, Seat: &models.Seat{SeatNumber: "A1"}},
		{ID: 2, CinemaID: 1, SeatID: 2, ShowDate: showDate, ShowTime: "19:00", IsAvailable: false, Seat: &models.Seat{SeatNumber: "A2"}},
	}

	repo.On("GetSeatAvailability", mock.Anything, 1, showDate, "19:00").Return(availabilities, nil)

	resp, err := service.GetSeatAvailability(context.Background(), 1, dateStr, "19:00")

	assert.NoError(t, err)
	assert.Equal(t, 1, resp.TotalAvailable)
	assert.Equal(t, 1, resp.TotalUnavailable)
	assert.Len(t, resp.AvailableSeats, 1)
	assert.Len(t, resp.UnavailableSeats, 1)
	repo.AssertExpectations(t)
}

func TestGetSeatAvailability_InvalidDate(t *testing.T) {
	repo := new(MockSeatAvailabilityRepository)
	service := NewSeatService(repo)

	resp, err := service.GetSeatAvailability(context.Background(), 1, "15-01-2026", "19:00")

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestGetSeatAvailability_RepoError(t *testing.T) {
	repo := new(MockSeatAvailabilityRepository)
	service := NewSeatService(repo)

	dateStr := "2026-01-15"
	showDate, _ := time.Parse("2006-01-02", dateStr)

	repo.On("GetSeatAvailability", mock.Anything, 1, showDate, "19:00").Return(nil, errors.New("db error"))

	resp, err := service.GetSeatAvailability(context.Background(), 1, dateStr, "19:00")

	assert.Error(t, err)
	assert.Nil(t, resp)
	repo.AssertExpectations(t)
}

func TestGetSeatByID_Success(t *testing.T) {
	repo := new(MockSeatAvailabilityRepository)
	service := NewSeatService(repo)

	seat := &models.Seat{ID: 1, SeatNumber: "A1"}
	repo.On("GetSeatByID", mock.Anything, 1).Return(seat, nil)

	result, err := service.GetSeatByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, seat, result)
	repo.AssertExpectations(t)
}

func TestGetSeatByID_Error(t *testing.T) {
	repo := new(MockSeatAvailabilityRepository)
	service := NewSeatService(repo)

	repo.On("GetSeatByID", mock.Anything, 1).Return(nil, errors.New("db error"))

	result, err := service.GetSeatByID(context.Background(), 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}
