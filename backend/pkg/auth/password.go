package auth

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	DefaultCost = 12
)

type PasswordService struct {
	cost int
}

func NewPasswordService() *PasswordService {
	return &PasswordService{
		cost: DefaultCost,
	}
}

func (p *PasswordService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), p.cost)
	return string(bytes), err
}

func (p *PasswordService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (p *PasswordService) ValidatePasswordStrength(password string) error {
	if len(password) < 6 {
		return &PasswordError{Message: "Password must be at least 6 characters long"}
	}

	// You can add more validation rules here
	// hasUpper := false
	// hasLower := false
	// hasNumber := false
	// hasSpecial := false

	// for _, char := range password {
	// 	switch {
	// 	case unicode.IsUpper(char):
	// 		hasUpper = true
	// 	case unicode.IsLower(char):
	// 		hasLower = true
	// 	case unicode.IsNumber(char):
	// 		hasNumber = true
	// 	case unicode.IsPunct(char) || unicode.IsSymbol(char):
	// 		hasSpecial = true
	// 	}
	// }

	return nil
}

type PasswordError struct {
	Message string
}

func (e *PasswordError) Error() string {
	return e.Message
}