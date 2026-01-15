package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/andre/project-app-bioskop-golang/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles user-related business logic
type UserService struct {
	userRepo     UserRepository
	emailService EmailSender
	jwtSecret    string
}

// UserRepository defines the persistence behavior needed by the user domain
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	CreateSession(ctx context.Context, session *models.UserSession) error
	GetSessionByToken(ctx context.Context, token string) (*models.UserSession, error)
	DeleteSession(ctx context.Context, token string) error
}

// EmailSender captures the OTP sending capability; concrete EmailService satisfies this.
type EmailSender interface {
	SendOTP(ctx context.Context, userID int, email, username string) error
}

// NewUserService creates a new UserService
func NewUserService(userRepo UserRepository, emailService EmailSender, jwtSecret string) *UserService {
	return &UserService{
		userRepo:     userRepo,
		emailService: emailService,
		jwtSecret:    jwtSecret,
	}
}

// RegisterUser registers a new user
func (s *UserService) RegisterUser(ctx context.Context, req *models.UserRegisterRequest) (*models.User, error) {
	// Check if username already exists
	existingUser, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	existingUser, err = s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing email: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Send OTP email automatically after registration
	if s.emailService != nil {
		err = s.emailService.SendOTP(ctx, user.ID, user.Email, user.Username)
		if err != nil {
			// Log error but don't fail registration
			fmt.Printf("Warning: Failed to send OTP email: %v\n", err)
		}
	}

	user.Password = "" // Don't return password
	return user, nil
}

// LoginUser authenticates a user and returns a token
func (s *UserService) LoginUser(ctx context.Context, req *models.UserLoginRequest) (*models.UserLoginResponse, error) {
	// Get user by username
	user, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate token
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Save session
	session := &models.UserSession{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	err = s.userRepo.CreateSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &models.UserLoginResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	}, nil
}

// LogoutUser logs out a user by deleting the session
func (s *UserService) LogoutUser(ctx context.Context, token string) error {
	err := s.userRepo.DeleteSession(ctx, token)
	if err != nil {
		return fmt.Errorf("failed to logout user: %w", err)
	}
	return nil
}

// VerifyToken verifies a JWT token and returns the user ID
func (s *UserService) VerifyToken(ctx context.Context, tokenString string) (int, error) {
	// First check if session exists
	session, err := s.userRepo.GetSessionByToken(ctx, tokenString)
	if err != nil {
		return 0, fmt.Errorf("failed to get session: %w", err)
	}
	if session == nil {
		return 0, errors.New("invalid session")
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		return 0, errors.New("session expired")
	}

	// Parse JWT token
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	// Extract user ID from token
	userID, err := token.Claims.GetSubject()
	if err != nil {
		return 0, fmt.Errorf("failed to get user id from token: %w", err)
	}

	var id int
	_, err = fmt.Sscanf(userID, "%d", &id)
	if err != nil {
		return 0, fmt.Errorf("invalid user id in token: %w", err)
	}

	return id, nil
}

// generateToken generates a JWT token
func (s *UserService) generateToken(userID int) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user != nil {
		user.Password = ""
	}
	return user, nil
}
