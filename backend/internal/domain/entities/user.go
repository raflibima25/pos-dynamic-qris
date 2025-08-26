package entities

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleCashier UserRole = "cashier"
)

type User struct {
	ID        string         `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Name      string         `json:"name" gorm:"not null"`
	Role      UserRole       `json:"role" gorm:"type:varchar(50);not null;check:role IN ('admin', 'cashier')"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Transactions []Transaction `json:"transactions,omitempty" gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return
}

func NewUser(email, name, password string, role UserRole) *User {
	return &User{
		ID:       uuid.New().String(),
		Email:    email,
		Name:     name,
		Password: password,
		Role:     role,
		IsActive: true,
	}
}

func (u *User) IsValidRole() bool {
	return u.Role == RoleAdmin || u.Role == RoleCashier
}

func (u *User) CanManageProducts() bool {
	return u.Role == RoleAdmin
}

func (u *User) CanProcessTransactions() bool {
	return u.Role == RoleAdmin || u.Role == RoleCashier
}
