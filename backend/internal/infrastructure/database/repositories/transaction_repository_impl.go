package repositories

import (
	"context"
	"errors"

	"qris-pos-backend/internal/domain/entities"
	"qris-pos-backend/internal/domain/repositories"

	"gorm.io/gorm"
)

type transactionRepositoryImpl struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) repositories.TransactionRepository {
	return &transactionRepositoryImpl{db: db}
}

func (r *transactionRepositoryImpl) Create(ctx context.Context, transaction *entities.Transaction) error {
	// Use transaction to ensure data consistency
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create the transaction first
		if err := tx.Omit("Items").Create(transaction).Error; err != nil {
			return err
		}

		// Create transaction items if any
		// Set TransactionID for each item and let database handle UUID generation
		if len(transaction.Items) > 0 {
			for i := range transaction.Items {
				transaction.Items[i].TransactionID = transaction.ID
			}
			if err := tx.Create(&transaction.Items).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *transactionRepositoryImpl) GetByID(ctx context.Context, id string) (*entities.Transaction, error) {
	var transaction entities.Transaction
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&transaction).Error

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *transactionRepositoryImpl) GetByIDWithDetails(ctx context.Context, id string) (*entities.Transaction, error) {
	var transaction entities.Transaction
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Category").
		Preload("Payment").
		Preload("QRCode").
		Where("id = ?", id).
		First(&transaction).Error

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *transactionRepositoryImpl) Update(ctx context.Context, transaction *entities.Transaction) error {
	return r.db.WithContext(ctx).Save(transaction).Error
}

func (r *transactionRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entities.Transaction{}, "id = ?", id).Error
}

func (r *transactionRepositoryImpl) List(ctx context.Context, filters repositories.TransactionFilters) ([]entities.Transaction, error) {
	var transactions []entities.Transaction
	query := r.db.WithContext(ctx).
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Preload("Payment")

	// Apply filters
	if filters.UserID != "" {
		query = query.Where("user_id = ?", filters.UserID)
	}

	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}

	if filters.DateFrom != nil {
		query = query.Where("created_at >= ?", *filters.DateFrom)
	}

	if filters.DateTo != nil {
		query = query.Where("created_at <= ?", *filters.DateTo)
	}

	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}

	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	err := query.Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepositoryImpl) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]entities.Transaction, error) {
	var transactions []entities.Transaction
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		Preload("Payment").
		Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&transactions).Error

	return transactions, err
}

func (r *transactionRepositoryImpl) GetByStatus(ctx context.Context, status entities.TransactionStatus, limit, offset int) ([]entities.Transaction, error) {
	var transactions []entities.Transaction
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Where("status = ?", status).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&transactions).Error

	return transactions, err
}

func (r *transactionRepositoryImpl) AddItem(ctx context.Context, item *entities.TransactionItem) error {
	// Check if item already exists for this transaction and product
	var existingItem entities.TransactionItem
	err := r.db.WithContext(ctx).
		Where("transaction_id = ? AND product_id = ?", item.TransactionID, item.ProductID).
		First(&existingItem).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err == nil {
		// Item exists, update quantity
		existingItem.Quantity += item.Quantity
		existingItem.TotalPrice = existingItem.UnitPrice * float64(existingItem.Quantity)
		return r.db.WithContext(ctx).Save(&existingItem).Error
	}

	// Item doesn't exist, create new
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *transactionRepositoryImpl) RemoveItem(ctx context.Context, transactionID, productID string) error {
	return r.db.WithContext(ctx).
		Where("transaction_id = ? AND product_id = ?", transactionID, productID).
		Delete(&entities.TransactionItem{}).Error
}

func (r *transactionRepositoryImpl) UpdateItemQuantity(ctx context.Context, transactionID, productID string, quantity int) error {
	if quantity <= 0 {
		return r.RemoveItem(ctx, transactionID, productID)
	}

	var item entities.TransactionItem
	err := r.db.WithContext(ctx).
		Where("transaction_id = ? AND product_id = ?", transactionID, productID).
		First(&item).Error

	if err != nil {
		return err
	}

	item.Quantity = quantity
	item.TotalPrice = item.UnitPrice * float64(quantity)

	return r.db.WithContext(ctx).Save(&item).Error
}

func (r *transactionRepositoryImpl) GetItems(ctx context.Context, transactionID string) ([]entities.TransactionItem, error) {
	var items []entities.TransactionItem
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Product.Category").
		Where("transaction_id = ?", transactionID).
		Find(&items).Error

	return items, err
}