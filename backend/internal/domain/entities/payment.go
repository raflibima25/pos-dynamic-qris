package entities

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentSuccess   PaymentStatus = "success"
	PaymentFailed    PaymentStatus = "failed"
	PaymentExpired   PaymentStatus = "expired"
	PaymentCancelled PaymentStatus = "cancelled"
)

type PaymentMethod string

const (
	PaymentMethodQRIS PaymentMethod = "qris"
)

type Payment struct {
	ID               string         `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TransactionID    string         `json:"transaction_id" gorm:"type:uuid;not null"`
	Amount           float64        `json:"amount" gorm:"type:decimal(10,2);not null;check:amount >= 0"`
	Method           PaymentMethod  `json:"method" gorm:"type:varchar(50);not null;check:method IN ('qris')"`
	Status           PaymentStatus  `json:"status" gorm:"type:varchar(50);not null;check:status IN ('pending', 'success', 'failed', 'expired', 'cancelled')"`
	ExternalID       string         `json:"external_id"`           // Midtrans transaction ID
	ExternalResponse string         `json:"external_response"`     // Midtrans response JSON
	PaidAt           *time.Time     `json:"paid_at"`
	ExpiresAt        time.Time      `json:"expires_at" gorm:"not null"`
	CreatedAt        time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Transaction Transaction `json:"transaction,omitempty" gorm:"foreignKey:TransactionID"`
	QRCode      *QRISCode   `json:"qr_code,omitempty" gorm:"foreignKey:PaymentID"`
}

func (Payment) TableName() string {
	return "payments"
}

func (p *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return
}

type QRISCode struct {
	ID            string         `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TransactionID string         `json:"transaction_id" gorm:"type:uuid;not null"`
	PaymentID     string         `json:"payment_id" gorm:"type:uuid;not null"`
	QRCode        string         `json:"qr_code" gorm:"not null"`
	QRImage       string         `json:"qr_image"`                                  // Base64 encoded image
	ExpiresAt     time.Time      `json:"expires_at" gorm:"not null"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Transaction Transaction `json:"transaction,omitempty" gorm:"foreignKey:TransactionID"`
	Payment     Payment     `json:"payment,omitempty" gorm:"foreignKey:PaymentID"`
}

func (QRISCode) TableName() string {
	return "qris_codes"
}

func (q *QRISCode) BeforeCreate(tx *gorm.DB) (err error) {
	if q.ID == "" {
		q.ID = uuid.New().String()
	}
	return
}

func NewPayment(transactionID string, amount float64, expiryMinutes int) *Payment {
	now := time.Now()
	expiresAt := now.Add(time.Duration(expiryMinutes) * time.Minute)

	return &Payment{
		ID:            uuid.New().String(),
		TransactionID: transactionID,
		Amount:        amount,
		Method:        PaymentMethodQRIS,
		Status:        PaymentPending,
		ExpiresAt:     expiresAt,
	}
}

func (p *Payment) IsExpired() bool {
	return time.Now().After(p.ExpiresAt)
}

func (p *Payment) CanBeProcessed() bool {
	return p.Status == PaymentPending && !p.IsExpired()
}

func (p *Payment) MarkAsSuccess(externalID, externalResponse string) {
	now := time.Now()
	p.Status = PaymentSuccess
	p.ExternalID = externalID
	p.ExternalResponse = externalResponse
	p.PaidAt = &now
}

func (p *Payment) MarkAsFailed(externalResponse string) {
	p.Status = PaymentFailed
	p.ExternalResponse = externalResponse
}

func (p *Payment) MarkAsExpired() {
	p.Status = PaymentExpired
}

func NewQRISCode(transactionID, paymentID, qrCode, qrImage string, expiryMinutes int) *QRISCode {
	now := time.Now()
	expiresAt := now.Add(time.Duration(expiryMinutes) * time.Minute)

	return &QRISCode{
		ID:            uuid.New().String(),
		TransactionID: transactionID,
		PaymentID:     paymentID,
		QRCode:        qrCode,
		QRImage:       qrImage,
		ExpiresAt:     expiresAt,
	}
}

func (q *QRISCode) IsExpired() bool {
	return time.Now().After(q.ExpiresAt)
}
