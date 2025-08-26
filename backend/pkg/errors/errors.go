package errors

import (
	"errors"
	"fmt"
)

var (
	// Authentication errors
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")

	// Authorization errors
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrInvalidRole    = errors.New("invalid role")

	// Validation errors
	ErrInvalidInput    = errors.New("invalid input")
	ErrRequiredField   = errors.New("required field missing")
	ErrInvalidFormat   = errors.New("invalid format")

	// Product errors
	ErrProductNotFound    = errors.New("product not found")
	ErrInsufficientStock  = errors.New("insufficient stock")
	ErrSKUExists          = errors.New("SKU already exists")

	// Transaction errors
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrEmptyCart           = errors.New("cart is empty")
	ErrTransactionExpired  = errors.New("transaction expired")

	// Payment errors
	ErrPaymentFailed   = errors.New("payment failed")
	ErrPaymentExpired  = errors.New("payment expired")
	ErrQRISExpired     = errors.New("QRIS code expired")
)

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code, message string, details any) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func NewValidationError(field, message string) *AppError {
	return &AppError{
		Code:    "VALIDATION_ERROR",
		Message: fmt.Sprintf("Validation failed for field '%s': %s", field, message),
		Details: map[string]string{"field": field, "error": message},
	}
}