package repositories

import (
	"context"
	"qris-pos-backend/internal/domain/entities"
	"qris-pos-backend/internal/domain/repositories"

	"gorm.io/gorm"
)

type paymentRepositoryImpl struct {
	db *gorm.DB
}

// NewPaymentRepository creates a new payment repository instance
func NewPaymentRepository(db *gorm.DB) repositories.PaymentRepository {
	return &paymentRepositoryImpl{db: db}
}

// CreatePayment creates a new payment record
func (r *paymentRepositoryImpl) CreatePayment(ctx context.Context, payment *entities.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

// GetPaymentByID retrieves a payment by its ID
func (r *paymentRepositoryImpl) GetPaymentByID(ctx context.Context, id string) (*entities.Payment, error) {
	var payment entities.Payment
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// GetPaymentByTransactionID retrieves a payment by transaction ID
func (r *paymentRepositoryImpl) GetPaymentByTransactionID(ctx context.Context, transactionID string) (*entities.Payment, error) {
	var payment entities.Payment
	err := r.db.WithContext(ctx).Where("transaction_id = ?", transactionID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// UpdatePayment updates a payment record
func (r *paymentRepositoryImpl) UpdatePayment(ctx context.Context, payment *entities.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}

// DeletePayment deletes a payment record
func (r *paymentRepositoryImpl) DeletePayment(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.Payment{}).Error
}

// CreateQRISCode creates a new QRIS code record
func (r *paymentRepositoryImpl) CreateQRISCode(ctx context.Context, qrisCode *entities.QRISCode) error {
	return r.db.WithContext(ctx).Create(qrisCode).Error
}

// GetQRISCodeByID retrieves a QRIS code by its ID
func (r *paymentRepositoryImpl) GetQRISCodeByID(ctx context.Context, id string) (*entities.QRISCode, error) {
	var qrisCode entities.QRISCode
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&qrisCode).Error
	if err != nil {
		return nil, err
	}
	return &qrisCode, nil
}

// GetQRISCodeByTransactionID retrieves a QRIS code by transaction ID
func (r *paymentRepositoryImpl) GetQRISCodeByTransactionID(ctx context.Context, transactionID string) (*entities.QRISCode, error) {
	var qrisCode entities.QRISCode
	err := r.db.WithContext(ctx).Where("transaction_id = ?", transactionID).First(&qrisCode).Error
	if err != nil {
		return nil, err
	}
	return &qrisCode, nil
}

// GetQRISCodeByPaymentID retrieves a QRIS code by payment ID
func (r *paymentRepositoryImpl) GetQRISCodeByPaymentID(ctx context.Context, paymentID string) (*entities.QRISCode, error) {
	var qrisCode entities.QRISCode
	err := r.db.WithContext(ctx).Where("payment_id = ?", paymentID).First(&qrisCode).Error
	if err != nil {
		return nil, err
	}
	return &qrisCode, nil
}

// UpdateQRISCode updates a QRIS code record
func (r *paymentRepositoryImpl) UpdateQRISCode(ctx context.Context, qrisCode *entities.QRISCode) error {
	return r.db.WithContext(ctx).Save(qrisCode).Error
}

// DeleteQRISCode deletes a QRIS code record
func (r *paymentRepositoryImpl) DeleteQRISCode(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.QRISCode{}).Error
}