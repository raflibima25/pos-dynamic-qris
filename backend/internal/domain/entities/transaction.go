package entities

import (
	"errors"
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type TransactionStatus string

const (
	StatusPending   TransactionStatus = "pending"
	StatusPaid      TransactionStatus = "paid" 
	StatusCancelled TransactionStatus = "cancelled"
	StatusExpired   TransactionStatus = "expired"
)

type Transaction struct {
	ID          string            `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      string            `json:"user_id" gorm:"type:uuid;not null"`
	TotalAmount float64           `json:"total_amount" gorm:"type:decimal(10,2);not null;check:total_amount >= 0"`
	TaxAmount   float64           `json:"tax_amount" gorm:"type:decimal(10,2);default:0;check:tax_amount >= 0"`
	Discount    float64           `json:"discount" gorm:"type:decimal(10,2);default:0;check:discount >= 0"`
	Status      TransactionStatus `json:"status" gorm:"type:varchar(50);not null;check:status IN ('pending', 'paid', 'cancelled', 'expired')"`
	Notes       string            `json:"notes"`
	CreatedAt   time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt    `json:"-" gorm:"index"`
	
	// Relations
	User     User              `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Items    []TransactionItem `json:"items,omitempty" gorm:"foreignKey:TransactionID"`
	Payment  *Payment          `json:"payment,omitempty" gorm:"foreignKey:TransactionID"`
	QRCode   *QRISCode         `json:"qr_code,omitempty" gorm:"foreignKey:TransactionID"`
}

func (Transaction) TableName() string {
	return "transactions"
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return
}

type TransactionItem struct {
	ID            string         `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TransactionID string         `json:"transaction_id" gorm:"type:uuid;not null"`
	ProductID     string         `json:"product_id" gorm:"type:uuid;not null"`
	Quantity      int            `json:"quantity" gorm:"not null;check:quantity > 0"`
	UnitPrice     float64        `json:"unit_price" gorm:"type:decimal(10,2);not null;check:unit_price >= 0"`
	TotalPrice    float64        `json:"total_price" gorm:"type:decimal(10,2);not null;check:total_price >= 0"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relations
	Transaction Transaction `json:"transaction,omitempty" gorm:"foreignKey:TransactionID"`
	Product     Product     `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

func (TransactionItem) TableName() string {
	return "transaction_items"
}

func (ti *TransactionItem) BeforeCreate(tx *gorm.DB) (err error) {
	if ti.ID == "" {
		ti.ID = uuid.New().String()
	}
	return
}

func NewTransaction(userID string) *Transaction {
	return &Transaction{
		ID:          uuid.New().String(),
		UserID:      userID,
		TotalAmount: 0,
		TaxAmount:   0,
		Discount:    0,
		Status:      StatusPending,
		Items:       []TransactionItem{},
	}
}

func (t *Transaction) AddItem(productID string, product *Product, quantity int) error {
	if product == nil {
		return errors.New("product cannot be nil")
	}
	
	if !product.IsAvailable() {
		return errors.New("product is not available")
	}
	
	if !product.CanFulfillQuantity(quantity) {
		return errors.New("insufficient stock")
	}
	
	unitPrice := product.Price
	totalPrice := unitPrice * float64(quantity)
	
	item := TransactionItem{
		ID:            uuid.New().String(),
		TransactionID: t.ID,
		ProductID:     productID,
		Quantity:      quantity,
		UnitPrice:     unitPrice,
		TotalPrice:    totalPrice,
		Product:       *product,
	}
	
	t.Items = append(t.Items, item)
	t.calculateTotal()
	
	return nil
}

func (t *Transaction) RemoveItem(productID string) {
	for i, item := range t.Items {
		if item.ProductID == productID {
			t.Items = append(t.Items[:i], t.Items[i+1:]...)
			break
		}
	}
	t.calculateTotal()
}

func (t *Transaction) calculateTotal() {
	var subtotal float64
	for _, item := range t.Items {
		subtotal += item.TotalPrice
	}
	
	t.TotalAmount = subtotal - t.Discount + t.TaxAmount
	t.UpdatedAt = time.Now()
}

func (t *Transaction) ApplyDiscount(discount float64) error {
	if discount < 0 {
		return errors.New("discount cannot be negative")
	}
	
	subtotal := t.getSubtotal()
	if discount > subtotal {
		return errors.New("discount cannot exceed subtotal")
	}
	
	t.Discount = discount
	t.calculateTotal()
	return nil
}

func (t *Transaction) ApplyTax(taxRate float64) error {
	if taxRate < 0 {
		return errors.New("tax rate cannot be negative")
	}
	
	subtotal := t.getSubtotal()
	t.TaxAmount = (subtotal - t.Discount) * taxRate / 100
	t.calculateTotal()
	return nil
}

func (t *Transaction) getSubtotal() float64 {
	var subtotal float64
	for _, item := range t.Items {
		subtotal += item.TotalPrice
	}
	return subtotal
}

func (t *Transaction) Cancel() error {
	if t.Status != StatusPending {
		return errors.New("only pending transactions can be cancelled")
	}
	
	t.Status = StatusCancelled
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Transaction) MarkAsPaid() error {
	if t.Status != StatusPending {
		return errors.New("only pending transactions can be marked as paid")
	}
	
	t.Status = StatusPaid
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Transaction) MarkAsExpired() error {
	if t.Status != StatusPending {
		return errors.New("only pending transactions can be marked as expired")
	}
	
	t.Status = StatusExpired
	t.UpdatedAt = time.Now()
	return nil
}