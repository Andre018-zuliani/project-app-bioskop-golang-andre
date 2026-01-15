package services

import (
	"context"
	"errors"
	"testing"

	"github.com/andre/project-app-bioskop-golang/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCinemaRepo struct {
	mock.Mock
}

func (m *MockCinemaRepo) GetAllCinemas(ctx context.Context, page, limit int, filters *models.CinemaFilters) ([]*models.Cinema, int, error) {
	args := m.Called(ctx, page, limit, filters)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*models.Cinema), args.Int(1), args.Error(2)
}

func (m *MockCinemaRepo) GetCinemaByID(ctx context.Context, id int) (*models.Cinema, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Cinema), args.Error(1)
}

func TestGetAllCinemas_Success(t *testing.T) {
	repo := new(MockCinemaRepo)
	service := NewCinemaService(repo)

	cinemas := []*models.Cinema{{ID: 1, Name: "Cinema A"}}
	filters := &models.CinemaFilters{City: "Jakarta"}
	repo.On("GetAllCinemas", mock.Anything, 1, 5, filters).Return(cinemas, 1, nil)

	resp, err := service.GetAllCinemas(context.Background(), 1, 5, filters)

	assert.NoError(t, err)
	assert.Equal(t, 1, resp.Total)
	assert.Equal(t, 1, resp.TotalPages)
	assert.Len(t, resp.Data.([]*models.Cinema), 1)
	repo.AssertExpectations(t)
}

func TestGetAllCinemas_Error(t *testing.T) {
	repo := new(MockCinemaRepo)
	service := NewCinemaService(repo)

	filters := &models.CinemaFilters{}
	repo.On("GetAllCinemas", mock.Anything, 1, 10, filters).Return(nil, 0, errors.New("db error"))

	resp, err := service.GetAllCinemas(context.Background(), 1, 10, filters)

	assert.Error(t, err)
	assert.Nil(t, resp)
	repo.AssertExpectations(t)
}

func TestGetCinemaByID_Success(t *testing.T) {
	repo := new(MockCinemaRepo)
	service := NewCinemaService(repo)

	cinema := &models.Cinema{ID: 2, Name: "Cinema B"}
	repo.On("GetCinemaByID", mock.Anything, 2).Return(cinema, nil)

	result, err := service.GetCinemaByID(context.Background(), 2)

	assert.NoError(t, err)
	assert.Equal(t, cinema, result)
	repo.AssertExpectations(t)
}

func TestGetCinemaByID_Error(t *testing.T) {
	repo := new(MockCinemaRepo)
	service := NewCinemaService(repo)

	repo.On("GetCinemaByID", mock.Anything, 2).Return(nil, errors.New("not found"))

	result, err := service.GetCinemaByID(context.Background(), 2)

	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}
