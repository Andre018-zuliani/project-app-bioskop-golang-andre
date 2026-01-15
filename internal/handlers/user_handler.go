package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/andre/project-app-bioskop-golang/internal/middleware"
	"github.com/andre/project-app-bioskop-golang/internal/models"
	"github.com/andre/project-app-bioskop-golang/internal/services"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService *services.UserService
	validator   *validator.Validate
	logger      *zap.Logger
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService *services.UserService, validator *validator.Validate, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator,
		logger:      logger,
	}
}

// Register handles user registration
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request", zap.Error(err))
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		h.logger.Error("validation error", zap.Error(err))
		writeError(w, "Validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Register user
	user, err := h.userService.RegisterUser(r.Context(), &req)
	if err != nil {
		h.logger.Error("failed to register user", zap.Error(err))
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Info("user registered successfully", zap.Int("user_id", user.ID))
	writeJSON(w, user, http.StatusCreated)
}

// Login handles user login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request", zap.Error(err))
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		h.logger.Error("validation error", zap.Error(err))
		writeError(w, "Validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Login user
	response, err := h.userService.LoginUser(r.Context(), &req)
	if err != nil {
		h.logger.Error("failed to login user", zap.Error(err))
		writeError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	h.logger.Info("user logged in successfully", zap.Int("user_id", response.ID))
	writeJSON(w, response, http.StatusOK)
}

// Logout handles user logout
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get token from context
	token, err := middleware.GetTokenFromContext(r)
	if err != nil {
		h.logger.Error("failed to get token from context", zap.Error(err))
		writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Logout user
	err = h.userService.LogoutUser(r.Context(), token)
	if err != nil {
		h.logger.Error("failed to logout user", zap.Error(err))
		writeError(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	h.logger.Info("user logged out successfully")
	writeJSON(w, map[string]string{"message": "Logout successful"}, http.StatusOK)
}

// GetProfile handles getting user profile
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(r)
	if err != nil {
		h.logger.Error("failed to get user id from context", zap.Error(err))
		writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user
	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		h.logger.Error("failed to get user", zap.Error(err))
		writeError(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	if user == nil {
		writeError(w, "User not found", http.StatusNotFound)
		return
	}

	writeJSON(w, user, http.StatusOK)
}
