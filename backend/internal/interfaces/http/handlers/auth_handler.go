package handlers

import (
	"qris-pos-backend/internal/interfaces/middleware"
	"qris-pos-backend/internal/usecases/auth"
	"qris-pos-backend/pkg/logger"
	"qris-pos-backend/pkg/response"
	"qris-pos-backend/pkg/validator"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUseCase *auth.AuthUseCase
	logger      logger.Logger
}

func NewAuthHandler(authUseCase *auth.AuthUseCase, logger logger.Logger) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
		logger:      logger,
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "Login request"
// @Success 200 {object} response.Response{data=auth.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind login request", "error", err)
		response.BadRequest(c, "Invalid request format", err.Error())
		return
	}

	// Validate request
	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.ValidationError(c, errors)
		return
	}

	result, err := h.authUseCase.Login(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Login failed", "error", err, "email", req.Email)
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, "Login successful", result)
}

// Register godoc
// @Summary User registration
// @Description Register a new user (Admin only)
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.RegisterRequest true "Registration request"
// @Success 201 {object} response.Response{data=auth.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req auth.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind register request", "error", err)
		response.BadRequest(c, "Invalid request format", err.Error())
		return
	}

	// Validate request
	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.ValidationError(c, errors)
		return
	}

	result, err := h.authUseCase.Register(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Registration failed", "error", err, "email", req.Email)
		if err.Error() == "email already exists" {
			response.BadRequest(c, "Email already exists", nil)
		} else {
			response.BadRequest(c, err.Error(), nil)
		}
		return
	}

	response.Created(c, "User registered successfully", result)
}

// GetProfile godoc
// @Summary Get current user profile
// @Description Get the profile of the currently authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=auth.UserResponse}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /auth/me [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	currentUser, exists := middleware.GetCurrentUser(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	result, err := h.authUseCase.GetCurrentUser(c.Request.Context(), currentUser.UserID)
	if err != nil {
		h.logger.Error("Failed to get user profile", "error", err, "user_id", currentUser.UserID)
		response.NotFound(c, "User not found")
		return
	}

	response.Success(c, "Profile retrieved successfully", result)
}

// RefreshToken godoc
// @Summary Refresh JWT token
// @Description Refresh the JWT token if it's close to expiry
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 401 {object} response.Response
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Get token from header
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		response.Unauthorized(c, "Invalid authorization header")
		return
	}

	token := authHeader[7:]
	newToken, err := h.authUseCase.RefreshToken(c.Request.Context(), token)
	if err != nil {
		h.logger.Error("Failed to refresh token", "error", err)
		response.Unauthorized(c, "Invalid token")
		return
	}

	response.Success(c, "Token refreshed successfully", map[string]string{
		"token": newToken,
	})
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=6"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change the password of the currently authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body ChangePasswordRequest true "Change password request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	currentUser, exists := middleware.GetCurrentUser(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format", err.Error())
		return
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.ValidationError(c, errors)
		return
	}

	if req.OldPassword == req.NewPassword {
		response.BadRequest(c, "New password must be different from old password", nil)
		return
	}

	err := h.authUseCase.ChangePassword(c.Request.Context(), currentUser.UserID, req.OldPassword, req.NewPassword)
	if err != nil {
		h.logger.Error("Failed to change password", "error", err, "user_id", currentUser.UserID)
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Password changed successfully", nil)
}

type UpdateProfileRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the profile of the currently authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body UpdateProfileRequest true "Update profile request"
// @Success 200 {object} response.Response{data=auth.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	currentUser, exists := middleware.GetCurrentUser(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format", err.Error())
		return
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.ValidationError(c, errors)
		return
	}

	result, err := h.authUseCase.UpdateProfile(c.Request.Context(), currentUser.UserID, req.Name)
	if err != nil {
		h.logger.Error("Failed to update profile", "error", err, "user_id", currentUser.UserID)
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Profile updated successfully", result)
}

// Logout godoc
// @Summary User logout
// @Description Logout user (client-side token removal)
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// In JWT, logout is typically handled on the client side
	// by removing the token from storage
	// For server-side logout, you would need to implement token blacklisting
	response.Success(c, "Logged out successfully", nil)
}