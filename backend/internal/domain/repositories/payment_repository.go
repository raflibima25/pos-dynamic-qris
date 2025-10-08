package repositories

import (
	"context"
	"qris-pos-backend/internal/domain/entities"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *entities.Payment) error
	GetPaymentByID(ctx context.Context, id string) (*entities.Payment, error)
	GetPaymentByTransactionID(ctx context.Context, transactionID string) (*entities.Payment, error)
	UpdatePayment(ctx context.Context, payment *entities.Payment) error
	DeletePayment(ctx context.Context, id string) error
	
	CreateQRISCode(ctx context.Context, qrisCode *entities.QRISCode) error
	GetQRISCodeByID(ctx context.Context, id string) (*entities.QRISCode, error)
	GetQRISCodeByTransactionID(ctx context.Context, transactionID string) (*entities.QRISCode, error)
	GetQRISCodeByPaymentID(ctx context.Context, paymentID string) (*entities.QRISCode, error)
	UpdateQRISCode(ctx context.Context, qrisCode *entities.QRISCode) error
	DeleteQRISCode(ctx context.Context, id string) error
}