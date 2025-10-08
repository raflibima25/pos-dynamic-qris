package transaction

import (
	"context"
	"errors"
	"fmt"

	"qris-pos-backend/internal/domain/entities"
	"qris-pos-backend/internal/domain/repositories"
	appErrors "qris-pos-backend/pkg/errors"
	"qris-pos-backend/pkg/logger"

	"gorm.io/gorm"
)

type CreateTransactionRequest struct {
	UserID string              `json:"user_id" validate:"required,uuid"`
	Items  []TransactionItemReq `json:"items" validate:"required,min=1"`
	Notes  string              `json:"notes"`
}

type TransactionItemReq struct {
	ProductID string `json:"product_id" validate:"required,uuid"`
	Quantity  int    `json:"quantity" validate:"required,gte=1"`
}

type AddItemRequest struct {
	ProductID string `json:"product_id" validate:"required,uuid"`
	Quantity  int    `json:"quantity" validate:"required,gte=1"`
}

type UpdateItemRequest struct {
	Quantity int `json:"quantity" validate:"required,gte=0"`
}

type TransactionResponse struct {
	ID          string                    `json:"id"`
	UserID      string                    `json:"user_id"`
	TotalAmount float64                   `json:"total_amount"`
	TaxAmount   float64                   `json:"tax_amount"`
	Discount    float64                   `json:"discount"`
	Status      entities.TransactionStatus `json:"status"`
	Notes       string                    `json:"notes"`
	CreatedAt   string                    `json:"created_at"`
	UpdatedAt   string                    `json:"updated_at"`
	Items       []TransactionItemResponse `json:"items"`
	User        *UserInfo                 `json:"user,omitempty"`
}

type TransactionItemResponse struct {
	ID         string      `json:"id"`
	ProductID  string      `json:"product_id"`
	Quantity   int         `json:"quantity"`
	UnitPrice  float64     `json:"unit_price"`
	TotalPrice float64     `json:"total_price"`
	Product    *ProductInfo `json:"product,omitempty"`
}

type UserInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type ProductInfo struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	CategoryName string `json:"category_name,omitempty"`
}

type TransactionUseCase struct {
	transactionRepo repositories.TransactionRepository
	productRepo     repositories.ProductRepository
	userRepo        repositories.UserRepository
	logger          logger.Logger
}

func NewTransactionUseCase(
	transactionRepo repositories.TransactionRepository,
	productRepo repositories.ProductRepository,
	userRepo repositories.UserRepository,
	logger logger.Logger,
) *TransactionUseCase {
	return &TransactionUseCase{
		transactionRepo: transactionRepo,
		productRepo:     productRepo,
		userRepo:        userRepo,
		logger:          logger,
	}
}

func (uc *TransactionUseCase) CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*TransactionResponse, error) {
	// Validate user exists
	_, err := uc.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrUserNotFound
		}
		return nil, err
	}

	// Create new transaction
	transaction := entities.NewTransaction(req.UserID)
	transaction.Notes = req.Notes

	// Add items and calculate total
	for _, itemReq := range req.Items {
		product, err := uc.productRepo.GetByID(ctx, itemReq.ProductID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("product with ID %s not found", itemReq.ProductID)
			}
			return nil, err
		}

		if err := transaction.AddItem(itemReq.ProductID, product, itemReq.Quantity); err != nil {
			return nil, err
		}
	}

	// Save transaction
	if err := uc.transactionRepo.Create(ctx, transaction); err != nil {
		uc.logger.Error("Failed to create transaction", "error", err, "user_id", req.UserID)
		return nil, err
	}

	uc.logger.Info("Transaction created successfully", "transaction_id", transaction.ID, "user_id", req.UserID)

	// Get full transaction with all relations (User, Items, Product)
	fullTransaction, err := uc.transactionRepo.GetByIDWithDetails(ctx, transaction.ID)
	if err != nil {
		return nil, err
	}

	return uc.mapTransactionToResponse(fullTransaction), nil
}

func (uc *TransactionUseCase) GetTransaction(ctx context.Context, id string) (*TransactionResponse, error) {
	// Get transaction with all details
	transaction, err := uc.transactionRepo.GetByIDWithDetails(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrTransactionNotFound
		}
		return nil, err
	}

	return uc.mapTransactionToResponse(transaction), nil
}

func (uc *TransactionUseCase) AddItemToTransaction(ctx context.Context, transactionID string, req *AddItemRequest) (*TransactionResponse, error) {
	// Get transaction
	transaction, err := uc.transactionRepo.GetByID(ctx, transactionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrTransactionNotFound
		}
		return nil, err
	}

	// Check if transaction can be modified
	if transaction.Status != entities.StatusPending {
		return nil, errors.New("cannot modify non-pending transaction")
	}

	// Get product
	product, err := uc.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrProductNotFound
		}
		return nil, err
	}

	// Create transaction item
	item := &entities.TransactionItem{
		TransactionID: transactionID,
		ProductID:     req.ProductID,
		Quantity:      req.Quantity,
		UnitPrice:     product.Price,
		TotalPrice:    product.Price * float64(req.Quantity),
		Product:       *product,
	}

	// Add item to transaction
	if err := uc.transactionRepo.AddItem(ctx, item); err != nil {
		return nil, err
	}

	// Recalculate transaction total
	if err := uc.recalculateTransaction(ctx, transactionID); err != nil {
		return nil, err
	}

	// Return updated transaction
	return uc.GetTransaction(ctx, transactionID)
}

func (uc *TransactionUseCase) RemoveItemFromTransaction(ctx context.Context, transactionID, productID string) (*TransactionResponse, error) {
	// Check transaction exists and is pending
	transaction, err := uc.transactionRepo.GetByID(ctx, transactionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrTransactionNotFound
		}
		return nil, err
	}

	if transaction.Status != entities.StatusPending {
		return nil, errors.New("cannot modify non-pending transaction")
	}

	// Remove item
	if err := uc.transactionRepo.RemoveItem(ctx, transactionID, productID); err != nil {
		return nil, err
	}

	// Recalculate transaction total
	if err := uc.recalculateTransaction(ctx, transactionID); err != nil {
		return nil, err
	}

	return uc.GetTransaction(ctx, transactionID)
}

func (uc *TransactionUseCase) UpdateItemQuantity(ctx context.Context, transactionID, productID string, req *UpdateItemRequest) (*TransactionResponse, error) {
	// Check transaction exists and is pending
	transaction, err := uc.transactionRepo.GetByID(ctx, transactionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrTransactionNotFound
		}
		return nil, err
	}

	if transaction.Status != entities.StatusPending {
		return nil, errors.New("cannot modify non-pending transaction")
	}

	// Update item quantity
	if err := uc.transactionRepo.UpdateItemQuantity(ctx, transactionID, productID, req.Quantity); err != nil {
		return nil, err
	}

	// Recalculate transaction total
	if err := uc.recalculateTransaction(ctx, transactionID); err != nil {
		return nil, err
	}

	return uc.GetTransaction(ctx, transactionID)
}

func (uc *TransactionUseCase) CancelTransaction(ctx context.Context, id string) error {
	transaction, err := uc.transactionRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return appErrors.ErrTransactionNotFound
		}
		return err
	}

	if err := transaction.Cancel(); err != nil {
		return err
	}

	if err := uc.transactionRepo.Update(ctx, transaction); err != nil {
		return err
	}

	uc.logger.Info("Transaction cancelled", "transaction_id", id)
	return nil
}

func (uc *TransactionUseCase) ListTransactions(ctx context.Context, filters repositories.TransactionFilters) ([]TransactionResponse, error) {
	transactions, err := uc.transactionRepo.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	responses := make([]TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		responses[i] = *uc.mapTransactionToResponse(&transaction)
	}

	return responses, nil
}

func (uc *TransactionUseCase) recalculateTransaction(ctx context.Context, transactionID string) error {
	// Get all items
	items, err := uc.transactionRepo.GetItems(ctx, transactionID)
	if err != nil {
		return err
	}

	// Calculate total
	var total float64
	for _, item := range items {
		total += item.TotalPrice
	}

	// Get transaction and update total
	transaction, err := uc.transactionRepo.GetByID(ctx, transactionID)
	if err != nil {
		return err
	}

	transaction.TotalAmount = total - transaction.Discount + transaction.TaxAmount

	return uc.transactionRepo.Update(ctx, transaction)
}

func (uc *TransactionUseCase) mapTransactionToResponse(transaction *entities.Transaction) *TransactionResponse {
	response := &TransactionResponse{
		ID:          transaction.ID,
		UserID:      transaction.UserID,
		TotalAmount: transaction.TotalAmount,
		TaxAmount:   transaction.TaxAmount,
		Discount:    transaction.Discount,
		Status:      transaction.Status,
		Notes:       transaction.Notes,
		CreatedAt:   transaction.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   transaction.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Items:       []TransactionItemResponse{},
	}

	// Map user info
	if transaction.User.ID != "" {
		response.User = &UserInfo{
			ID:   transaction.User.ID,
			Name: transaction.User.Name,
			Role: string(transaction.User.Role),
		}
	}

	// Map items
	for _, item := range transaction.Items {
		itemResponse := TransactionItemResponse{
			ID:         item.ID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
			TotalPrice: item.TotalPrice,
		}

		// Map product info
		if item.Product.ID != "" {
			itemResponse.Product = &ProductInfo{
				ID:    item.Product.ID,
				Name:  item.Product.Name,
				Price: item.Product.Price,
				Stock: item.Product.Stock,
			}

			if item.Product.Category.ID != "" {
				itemResponse.Product.CategoryName = item.Product.Category.Name
			}
		}

		response.Items = append(response.Items, itemResponse)
	}

	return response
}