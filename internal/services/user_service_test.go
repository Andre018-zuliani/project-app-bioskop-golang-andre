package services

import (
	"context"
	"testing"
	"time"

	"github.com/andre/project-app-bioskop-golang/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) CreateSession(ctx context.Context, session *models.UserSession) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockUserRepository) GetSessionByToken(ctx context.Context, token string) (*models.UserSession, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserSession), args.Error(1)
}

func (m *MockUserRepository) DeleteSession(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func TestRegisterUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, nil, "test-secret")

	req := &models.UserRegisterRequest{Username: "testuser", Email: "test@example.com", Password: "password123"}

	mockRepo.On("GetUserByUsername", mock.Anything, "testuser").Return(nil, nil)
	mockRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(nil, nil)
	mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).Run(func(args mock.Arguments) {
		u := args.Get(1).(*models.User)
		u.ID = 1
	}).Return(nil)

	user, err := service.RegisterUser(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_UsernameExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, nil, "test-secret")

	req := &models.UserRegisterRequest{Username: "existinguser", Email: "test@example.com", Password: "password123"}
	existingUser := &models.User{ID: 1, Username: "existinguser"}

	mockRepo.On("GetUserByUsername", mock.Anything, "existinguser").Return(existingUser, nil)

	user, err := service.RegisterUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "username already exists")
	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_EmailExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, nil, "test-secret")

	req := &models.UserRegisterRequest{Username: "testuser", Email: "existing@example.com", Password: "password123"}
	existingUser := &models.User{ID: 1, Email: "existing@example.com"}

	mockRepo.On("GetUserByUsername", mock.Anything, "testuser").Return(nil, nil)
	mockRepo.On("GetUserByEmail", mock.Anything, "existing@example.com").Return(existingUser, nil)

	user, err := service.RegisterUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "email already exists")
	mockRepo.AssertExpectations(t)
}

func TestLoginUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, nil, "test-secret")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	existingUser := &models.User{ID: 1, Username: "testuser", Email: "test@example.com", Password: string(hashedPassword)}

	req := &models.UserLoginRequest{Username: "testuser", Password: "password123"}

	mockRepo.On("GetUserByUsername", mock.Anything, "testuser").Return(existingUser, nil)
	mockRepo.On("CreateSession", mock.Anything, mock.AnythingOfType("*models.UserSession")).Return(nil)

	response, err := service.LoginUser(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, "testuser", response.Username)
	mockRepo.AssertExpectations(t)
}

func TestLoginUser_InvalidCredentials(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, nil, "test-secret")

	req := &models.UserLoginRequest{Username: "nonexistent", Password: "password123"}

	mockRepo.On("GetUserByUsername", mock.Anything, "nonexistent").Return(nil, nil)

	response, err := service.LoginUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "invalid credentials")
	mockRepo.AssertExpectations(t)
}

func TestLoginUser_WrongPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, nil, "test-secret")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	existingUser := &models.User{ID: 1, Username: "testuser", Password: string(hashedPassword)}

	req := &models.UserLoginRequest{Username: "testuser", Password: "wrongpassword"}

	mockRepo.On("GetUserByUsername", mock.Anything, "testuser").Return(existingUser, nil)

	response, err := service.LoginUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "invalid credentials")
	mockRepo.AssertExpectations(t)
}

func TestVerifyToken_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, nil, "test-secret")

	token, _ := service.generateToken(1)
	session := &models.UserSession{UserID: 1, Token: token, ExpiresAt: time.Now().Add(24 * time.Hour)}

	mockRepo.On("GetSessionByToken", mock.Anything, token).Return(session, nil)

	userID, err := service.VerifyToken(context.Background(), token)

	assert.NoError(t, err)
	assert.Equal(t, 1, userID)
	mockRepo.AssertExpectations(t)
}

func TestVerifyToken_InvalidToken(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, nil, "test-secret")

	invalidToken := "invalid.token"
	mockRepo.On("GetSessionByToken", mock.Anything, invalidToken).Return(nil, nil)

	userID, err := service.VerifyToken(context.Background(), invalidToken)

	assert.Error(t, err)
	assert.Equal(t, 0, userID)
	mockRepo.AssertExpectations(t)
}

func TestLogoutUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, nil, "test-secret")

	token := "valid.token"
	mockRepo.On("DeleteSession", mock.Anything, token).Return(nil)

	err := service.LogoutUser(context.Background(), token)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, nil, "test-secret")

	expectedUser := &models.User{ID: 1, Username: "testuser", Email: "test@example.com", Password: "hidden"}
	mockRepo.On("GetUserByID", mock.Anything, 1).Return(expectedUser, nil)

	user, err := service.GetUserByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
	assert.Empty(t, user.Password)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, nil, "test-secret")

	mockRepo.On("GetUserByID", mock.Anything, 999).Return(nil, nil)

	user, err := service.GetUserByID(context.Background(), 999)

	assert.NoError(t, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}
