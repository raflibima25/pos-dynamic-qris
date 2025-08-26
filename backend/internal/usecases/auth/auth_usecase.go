package auth

import (
	"context"
	"errors"

	"qris-pos-backend/internal/domain/entities"
	"qris-pos-backend/internal/domain/repositories"
	"qris-pos-backend/pkg/auth"
	appErrors "qris-pos-backend/pkg/errors"
	"qris-pos-backend/pkg/logger"

	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Name     string            `json:"name" validate:"required,min=2,max=100"`
	Email    string            `json:"email" validate:"required,email"`
	Password string            `json:"password" validate:"required,min=6"`
	Role     entities.UserRole `json:"role" validate:"required,oneof=admin cashier"`
}

type LoginResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}

type UserResponse struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Email    string            `json:"email"`
	Role     entities.UserRole `json:"role"`
	IsActive bool              `json:"is_active"`
}

type AuthUseCase struct {
	userRepo        repositories.UserRepository
	passwordService *auth.PasswordService
	jwtService      *auth.JWTService
	logger          logger.Logger
}

func NewAuthUseCase(
	userRepo repositories.UserRepository,
	passwordService *auth.PasswordService,
	jwtService *auth.JWTService,
	logger logger.Logger,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:        userRepo,
		passwordService: passwordService,
		jwtService:      jwtService,
		logger:          logger,
	}
}

func (uc *AuthUseCase) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// Find user by email
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			uc.logger.Warn("Login attempt with non-existent email", "email", req.Email)
			return nil, appErrors.ErrInvalidCredentials
		}
		uc.logger.Error("Failed to get user by email", "error", err)
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		uc.logger.Warn("Login attempt with inactive user", "user_id", user.ID)
		return nil, appErrors.ErrInvalidCredentials
	}

	// Verify password
	if !uc.passwordService.CheckPasswordHash(req.Password, user.Password) {
		uc.logger.Warn("Invalid password attempt", "user_id", user.ID)
		return nil, appErrors.ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := uc.jwtService.GenerateToken(user)
	if err != nil {
		uc.logger.Error("Failed to generate JWT token", "error", err, "user_id", user.ID)
		return nil, errors.New("failed to generate token")
	}

	uc.logger.Info("User logged in successfully", "user_id", user.ID, "email", user.Email)

	return &LoginResponse{
		User:  uc.mapUserToResponse(user),
		Token: token,
	}, nil
}

func (uc *AuthUseCase) Register(ctx context.Context, req *RegisterRequest) (*UserResponse, error) {
	// Check if email already exists
	existingUser, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		uc.logger.Error("Failed to check existing user", "error", err)
		return nil, err
	}

	if existingUser != nil {
		uc.logger.Warn("Registration attempt with existing email", "email", req.Email)
		return nil, appErrors.ErrEmailExists
	}

	// Validate password strength
	if err := uc.passwordService.ValidatePasswordStrength(req.Password); err != nil {
		return nil, err
	}

	// Hash password
	hashedPassword, err := uc.passwordService.HashPassword(req.Password)
	if err != nil {
		uc.logger.Error("Failed to hash password", "error", err)
		return nil, errors.New("failed to process password")
	}

	// Create user
	user := entities.NewUser(req.Email, req.Name, hashedPassword, req.Role)

	if err := uc.userRepo.Create(ctx, user); err != nil {
		uc.logger.Error("Failed to create user", "error", err)
		return nil, err
	}

	uc.logger.Info("User registered successfully", "user_id", user.ID, "email", user.Email)

	return uc.mapUserToResponse(user), nil
}

func (uc *AuthUseCase) GetCurrentUser(ctx context.Context, userID string) (*UserResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrUserNotFound
		}
		uc.logger.Error("Failed to get user", "error", err, "user_id", userID)
		return nil, err
	}

	return uc.mapUserToResponse(user), nil
}

func (uc *AuthUseCase) RefreshToken(ctx context.Context, token string) (string, error) {
	newToken, err := uc.jwtService.RefreshToken(token)
	if err != nil {
		uc.logger.Error("Failed to refresh token", "error", err)
		return "", appErrors.ErrInvalidToken
	}

	return newToken, nil
}

func (uc *AuthUseCase) ChangePassword(ctx context.Context, userID string, oldPassword, newPassword string) error {
	// Get user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return appErrors.ErrUserNotFound
		}
		return err
	}

	// Verify old password
	if !uc.passwordService.CheckPasswordHash(oldPassword, user.Password) {
		return appErrors.ErrInvalidCredentials
	}

	// Validate new password
	if err := uc.passwordService.ValidatePasswordStrength(newPassword); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := uc.passwordService.HashPassword(newPassword)
	if err != nil {
		uc.logger.Error("Failed to hash new password", "error", err)
		return errors.New("failed to process password")
	}

	// Update password
	user.Password = hashedPassword
	if err := uc.userRepo.Update(ctx, user); err != nil {
		uc.logger.Error("Failed to update user password", "error", err)
		return err
	}

	uc.logger.Info("Password changed successfully", "user_id", userID)
	return nil
}

func (uc *AuthUseCase) UpdateProfile(ctx context.Context, userID string, name string) (*UserResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrUserNotFound
		}
		return nil, err
	}

	user.Name = name
	if err := uc.userRepo.Update(ctx, user); err != nil {
		uc.logger.Error("Failed to update user profile", "error", err)
		return nil, err
	}

	uc.logger.Info("Profile updated successfully", "user_id", userID)
	return uc.mapUserToResponse(user), nil
}

func (uc *AuthUseCase) mapUserToResponse(user *entities.User) *UserResponse {
	return &UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Role:     user.Role,
		IsActive: user.IsActive,
	}
}
