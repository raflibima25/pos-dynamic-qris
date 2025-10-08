package repositories

import (
	"context"
	"qris-pos-backend/internal/domain/entities"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *entities.Transaction) error
	GetByID(ctx context.Context, id string) (*entities.Transaction, error)
	GetByIDWithDetails(ctx context.Context, id string) (*entities.Transaction, error)
	Update(ctx context.Context, transaction *entities.Transaction) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filters TransactionFilters) ([]entities.Transaction, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]entities.Transaction, error)
	GetByStatus(ctx context.Context, status entities.TransactionStatus, limit, offset int) ([]entities.Transaction, error)

	// Transaction Items operations
	AddItem(ctx context.Context, item *entities.TransactionItem) error
	RemoveItem(ctx context.Context, transactionID, productID string) error
	UpdateItemQuantity(ctx context.Context, transactionID, productID string, quantity int) error
	GetItems(ctx context.Context, transactionID string) ([]entities.TransactionItem, error)
}

type TransactionFilters struct {
	UserID    string
	Status    entities.TransactionStatus
	DateFrom  *string // Format: "2023-01-01"
	DateTo    *string // Format: "2023-12-31"
	Limit     int
	Offset    int
}