package entities

import (
	"errors"
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type Product struct {
	ID          string         `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Price       float64        `json:"price" gorm:"type:decimal(10,2);not null;check:price >= 0"`
	Stock       int            `json:"stock" gorm:"not null;check:stock >= 0"`
	CategoryID  string         `json:"category_id" gorm:"type:uuid;not null"`
	SKU         string         `json:"sku" gorm:"uniqueIndex"`
	ImageURL    string         `json:"image_url" gorm:"type:text"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Category         Category          `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	TransactionItems []TransactionItem `json:"transaction_items,omitempty" gorm:"foreignKey:ProductID"`
}

func (Product) TableName() string {
	return "products"
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return
}

func NewProduct(name, description, sku, categoryID string, price float64, stock int) (*Product, error) {
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}
	if price < 0 {
		return nil, errors.New("product price cannot be negative")
	}
	if stock < 0 {
		return nil, errors.New("product stock cannot be negative")
	}

	return &Product{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		CategoryID:  categoryID,
		SKU:         sku,
		IsActive:    true,
	}, nil
}

func (p *Product) UpdateStock(quantity int) error {
	newStock := p.Stock + quantity
	if newStock < 0 {
		return errors.New("insufficient stock")
	}
	
	p.Stock = newStock
	return nil
}

func (p *Product) IsAvailable() bool {
	return p.IsActive && p.Stock > 0
}

func (p *Product) CanFulfillQuantity(quantity int) bool {
	return p.Stock >= quantity
}

type Category struct {
	ID        string         `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string         `json:"name" gorm:"uniqueIndex;not null"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Products []Product `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
}

func (Category) TableName() string {
	return "categories"
}

func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return
}