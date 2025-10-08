package payment

import (
	"context"
	"fmt"
	"qris-pos-backend/internal/domain/entities"
	"qris-pos-backend/internal/domain/repositories"
	"qris-pos-backend/internal/infrastructure/payment"
	"qris-pos-backend/internal/infrastructure/qrcode"
	appErrors "qris-pos-backend/pkg/errors"
	"qris-pos-backend/pkg/logger"
	"strings"
	"time"

	"gorm.io/gorm"
)

type GenerateQRISRequest struct {
	TransactionID string  `json:"transaction_id" validate:"required,uuid"`
	Amount        float64 `json:"amount" validate:"required,gte=0"`
	CallbackURL   string  `json:"callback_url"`
	ExpiryMinutes int     `json:"expiry_minutes"`
}

type PaymentResponse struct {
	ID            string                 `json:"id"`
	TransactionID string                 `json:"transaction_id"`
	Amount        float64                `json:"amount"`
	Method        entities.PaymentMethod `json:"method"`
	Status        entities.PaymentStatus `json:"status"`
	ExternalID    string                 `json:"external_id"`
	PaidAt        *string                `json:"paid_at"`
	ExpiresAt     string                 `json:"expires_at"`
	CreatedAt     string                 `json:"created_at"`
	UpdatedAt     string                 `json:"updated_at"`
	QRISCode      *QRISCodeResponse      `json:"qr_code,omitempty"`
}

type QRISCodeResponse struct {
	ID            string `json:"id"`
	TransactionID string `json:"transaction_id"`
	PaymentID     string `json:"payment_id"`
	QRCode        string `json:"qr_code"` // QRIS EMVCo string for frontend QR generation
	URL           string `json:"url"`     // Midtrans simulator URL for testing
	ExpiresAt     string `json:"expires_at"`
	CreatedAt     string `json:"created_at"`
}

type PaymentStatusResponse struct {
	TransactionID string                 `json:"transaction_id"`
	Status        entities.PaymentStatus `json:"status"`
	ExternalID    string                 `json:"external_id"`
	Message       string                 `json:"message"`
}

type PaymentUseCase struct {
	paymentRepo      repositories.PaymentRepository
	transactionRepo  repositories.TransactionRepository
	midtransClient   *payment.MidtransClient
	qrCodeGenerator  *qrcode.QRCodeGenerator
	logger           logger.Logger
	defaultExpiryMin int
}

func NewPaymentUseCase(
	paymentRepo repositories.PaymentRepository,
	transactionRepo repositories.TransactionRepository,
	midtransClient *payment.MidtransClient,
	qrCodeGenerator *qrcode.QRCodeGenerator,
	logger logger.Logger,
) *PaymentUseCase {
	return &PaymentUseCase{
		paymentRepo:      paymentRepo,
		transactionRepo:  transactionRepo,
		midtransClient:   midtransClient,
		qrCodeGenerator:  qrCodeGenerator,
		logger:           logger,
		defaultExpiryMin: 10, // Default 10 minutes expiry
	}
}

// GenerateQRIS generates a QRIS code for a transaction
func (uc *PaymentUseCase) GenerateQRIS(ctx context.Context, req *GenerateQRISRequest) (*PaymentResponse, error) {
	// Validate transaction exists and is pending
	// Use GetByIDWithDetails to preload User and Items
	transaction, err := uc.transactionRepo.GetByIDWithDetails(ctx, req.TransactionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, appErrors.ErrTransactionNotFound
		}
		return nil, err
	}

	// Check if transaction is in pending status
	if transaction.Status != entities.StatusPending {
		return nil, fmt.Errorf("transaction is not in pending status")
	}

	// Check if transaction already has a payment
	existingPayment, err := uc.paymentRepo.GetPaymentByTransactionID(ctx, req.TransactionID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if existingPayment != nil {
		// If payment exists and is still valid, return it
		if existingPayment.CanBeProcessed() {
			// Get existing QRIS code
			existingQRIS, err := uc.paymentRepo.GetQRISCodeByPaymentID(ctx, existingPayment.ID)
			if err != nil && err != gorm.ErrRecordNotFound {
				return nil, err
			}

			return uc.mapPaymentToResponse(existingPayment, existingQRIS), nil
		}

		// If payment is expired, mark it as expired
		if existingPayment.IsExpired() {
			existingPayment.MarkAsExpired()
			if err := uc.paymentRepo.UpdatePayment(ctx, existingPayment); err != nil {
				uc.logger.Error("Failed to update expired payment", "error", err)
			}
		}
	}

	// Determine expiry minutes
	expiryMinutes := req.ExpiryMinutes
	if expiryMinutes <= 0 {
		expiryMinutes = uc.defaultExpiryMin
	}

	// Create payment record
	paymentEntity := entities.NewPayment(req.TransactionID, req.Amount, expiryMinutes)

	// Generate QRIS via Midtrans
	// OrderID must be <= 50 chars. Using first 8 chars of UUID + current timestamp
	shortTxID := req.TransactionID
	if len(shortTxID) > 8 {
		shortTxID = shortTxID[:8]
	}
	// Use current time for consistent order_id
	now := time.Now()
	orderID := fmt.Sprintf("qris-%s-%d", shortTxID, now.Unix())

	// Store order_id in payment entity for later status checking
	paymentEntity.OrderID = orderID

	qrisReq := payment.QRISRequest{
		TransactionID:  req.TransactionID,
		OrderID:        orderID,
		GrossAmount:    transaction.TotalAmount, // Use transaction total (includes tax & discount)
		CustomerName:   transaction.User.Name,
		CustomerEmail:  transaction.User.Email,
		Items:          uc.mapTransactionItemsToQRISItems(transaction),
		ExpiryDuration: expiryMinutes,
	}

	// Debug: Log QRIS request details
	uc.logger.Info("Generating QRIS with details",
		"order_id", orderID,
		"gross_amount", qrisReq.GrossAmount,
		"items_count", len(qrisReq.Items),
		"transaction_total", transaction.TotalAmount)

	// Debug: Log each item
	var itemsSum float64
	for _, item := range qrisReq.Items {
		itemTotal := item.Price * float64(item.Quantity)
		itemsSum += itemTotal
		uc.logger.Info("Item details",
			"name", item.Name,
			"unit_price", item.Price,
			"quantity", item.Quantity,
			"item_total", itemTotal)
	}
	uc.logger.Info("Items sum check",
		"items_sum", itemsSum,
		"gross_amount", qrisReq.GrossAmount,
		"match", itemsSum == qrisReq.GrossAmount)

	qrisResponse, err := uc.midtransClient.GenerateQRIS(ctx, qrisReq)
	if err != nil {
		uc.logger.Error("Failed to generate QRIS via Midtrans", "error", err)
		return nil, fmt.Errorf("failed to generate QRIS: %w", err)
	}

	// Save payment first to get the ID
	if err := uc.paymentRepo.CreatePayment(ctx, paymentEntity); err != nil {
		// Check if error is due to duplicate constraint violation
		if strings.Contains(err.Error(), "idx_unique_pending_payment_per_transaction") {
			// Payment already exists, fetch and return existing payment
			uc.logger.Warn("Payment already exists for transaction, returning existing", "transaction_id", req.TransactionID)
			existingPayment, getErr := uc.paymentRepo.GetPaymentByTransactionID(ctx, req.TransactionID)
			if getErr != nil {
				uc.logger.Error("Failed to get existing payment", "error", getErr)
				return nil, err
			}
			existingQRIS, getErr := uc.paymentRepo.GetQRISCodeByPaymentID(ctx, existingPayment.ID)
			if getErr != nil {
				uc.logger.Error("Failed to get existing QRIS", "error", getErr)
				return nil, err
			}
			return uc.mapPaymentToResponse(existingPayment, existingQRIS), nil
		}
		uc.logger.Error("Failed to create payment record", "error", err)
		return nil, err
	}

	// Now create QRIS code record with the payment ID
	// Store both QRIS string (for frontend QR generation) and URL (for Midtrans simulator testing)
	qrCodeEntity := entities.NewQRISCode(
		req.TransactionID,
		paymentEntity.ID,
		qrisResponse.QRString,
		qrisResponse.URL, // Midtrans simulator URL for testing
		expiryMinutes,
	)

	// Save QRIS code
	if err := uc.paymentRepo.CreateQRISCode(ctx, qrCodeEntity); err != nil {
		uc.logger.Error("Failed to create QRIS code record", "error", err)
		// Try to rollback payment creation
		if delErr := uc.paymentRepo.DeletePayment(ctx, paymentEntity.ID); delErr != nil {
			uc.logger.Error("Failed to rollback payment creation", "error", delErr)
		}
		return nil, err
	}

	uc.logger.Info("QRIS generated successfully", "transaction_id", req.TransactionID, "payment_id", paymentEntity.ID)

	return uc.mapPaymentToResponse(paymentEntity, qrCodeEntity), nil
}

// GetPaymentStatus gets the status of a payment for a transaction
func (uc *PaymentUseCase) GetPaymentStatus(ctx context.Context, transactionID string) (*PaymentStatusResponse, error) {
	// Get payment record
	paymentEntity, err := uc.paymentRepo.GetPaymentByTransactionID(ctx, transactionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, appErrors.ErrPaymentNotFound
		}
		return nil, err
	}

	// If payment is already success, failed, or cancelled, return current status
	if paymentEntity.Status != entities.PaymentPending && paymentEntity.Status != entities.PaymentExpired {
		return &PaymentStatusResponse{
			TransactionID: transactionID,
			Status:        paymentEntity.Status,
			ExternalID:    paymentEntity.ExternalID,
			Message:       fmt.Sprintf("Payment status: %s", paymentEntity.Status),
		}, nil
	}

	// If payment is expired, return expired status
	if paymentEntity.IsExpired() {
		// Update payment status to expired if not already marked
		if paymentEntity.Status != entities.PaymentExpired {
			paymentEntity.MarkAsExpired()
			if err := uc.paymentRepo.UpdatePayment(ctx, paymentEntity); err != nil {
				uc.logger.Error("Failed to update expired payment", "error", err)
			}
		}

		return &PaymentStatusResponse{
			TransactionID: transactionID,
			Status:        entities.PaymentExpired,
			Message:       "Payment has expired",
		}, nil
	}

	// For pending payments, check status with Midtrans
	// Use the stored OrderID from payment entity
	orderID := paymentEntity.OrderID
	if orderID == "" {
		// Fallback: If OrderID is not stored (old data), return pending
		uc.logger.Warn("OrderID not found in payment, cannot check status", "transaction_id", transactionID)
		return &PaymentStatusResponse{
			TransactionID: transactionID,
			Status:        entities.PaymentPending,
			Message:       "Payment is pending. Waiting for customer to complete payment.",
		}, nil
	}

	// Check status with Midtrans
	midtransStatus, err := uc.midtransClient.GetTransactionStatus(ctx, orderID)
	if err != nil {
		uc.logger.Error("Failed to check Midtrans status", "error", err, "order_id", orderID)
		return &PaymentStatusResponse{
			TransactionID: transactionID,
			Status:        entities.PaymentPending,
			Message:       "Payment is pending. Waiting for customer to complete payment.",
		}, nil
	}

	// Update payment based on Midtrans status
	var newStatus entities.PaymentStatus
	switch midtransStatus.TransactionStatus {
	case "settlement", "capture":
		newStatus = entities.PaymentSuccess
		paymentEntity.MarkAsSuccess(midtransStatus.TransactionID, midtransStatus.StatusMessage)

		// Update transaction status
		transaction, _ := uc.transactionRepo.GetByID(ctx, transactionID)
		if transaction != nil {
			transaction.MarkAsPaid()
			uc.transactionRepo.Update(ctx, transaction)
		}
	case "pending":
		newStatus = entities.PaymentPending
	case "deny", "cancel", "expire":
		newStatus = entities.PaymentFailed
		paymentEntity.MarkAsFailed(midtransStatus.StatusMessage)
	default:
		newStatus = entities.PaymentPending
	}

	// Update payment in database
	if err := uc.paymentRepo.UpdatePayment(ctx, paymentEntity); err != nil {
		uc.logger.Error("Failed to update payment status", "error", err)
	}

	return &PaymentStatusResponse{
		TransactionID: transactionID,
		Status:        newStatus,
		ExternalID:    midtransStatus.TransactionID,
		Message:       midtransStatus.StatusMessage,
	}, nil
}

// HandlePaymentNotification handles payment notifications from Midtrans
func (uc *PaymentUseCase) HandlePaymentNotification(ctx context.Context, orderID string, status string, externalID string, response string) error {
	// Since we shortened the order_id, we need to find payment by external_id (Midtrans transaction_id)
	// which should be stored in the payment record
	// For simplicity, we'll use the externalID parameter to find the payment
	uc.logger.Info("Received payment notification", "order_id", orderID, "external_id", externalID, "status", status)

	// Find payment by external_id - this would require adding a new method to repository
	// For now, we'll need to extract what we can and handle accordingly
	// The externalID from Midtrans should help us identify the payment

	// For now, we'll use order_id to look up the QRIS code
	// Order ID format: qris-{short_tx_id}-{timestamp}
	// We can search QRIS codes by matching the order_id pattern

	// Temporary solution: Get all recent QRIS codes and match by external_id
	// This is not optimal but works for MVP
	// TODO: Add proper index and lookup by external_id or order_id

	// For now, just log and return - webhook implementation can be done later
	uc.logger.Warn("Payment notification received but lookup not fully implemented",
		"order_id", orderID,
		"external_id", externalID,
		"status", status)

	// Return nil to acknowledge receipt
	uc.logger.Info("Payment notification acknowledged", "order_id", orderID, "status", status)
	return nil
}

// RefreshQRIS refreshes an expired QRIS code
func (uc *PaymentUseCase) RefreshQRIS(ctx context.Context, transactionID string) (*PaymentResponse, error) {
	// Get existing payment
	paymentEntity, err := uc.paymentRepo.GetPaymentByTransactionID(ctx, transactionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, appErrors.ErrPaymentNotFound
		}
		return nil, err
	}

	// Check if payment can be refreshed (expired or about to expire)
	if paymentEntity.Status != entities.PaymentPending && paymentEntity.Status != entities.PaymentExpired {
		return nil, fmt.Errorf("payment cannot be refreshed: current status is %s", paymentEntity.Status)
	}

	// Get transaction with details for User and Items
	transaction, err := uc.transactionRepo.GetByIDWithDetails(ctx, transactionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, appErrors.ErrTransactionNotFound
		}
		return nil, err
	}

	// If payment is not expired yet, check if it's close to expiry (within 2 minutes)
	if paymentEntity.Status == entities.PaymentPending && !paymentEntity.IsExpired() {
		// For simplicity, we'll allow refresh anytime - in production you might want to restrict this
		// based on time to expiry
	}

	// Generate new QRIS via Midtrans
	// Use short transaction ID (first 8 chars) to keep order_id under 50 chars limit
	shortTxID := transactionID
	if len(shortTxID) > 8 {
		shortTxID = shortTxID[:8]
	}
	now := time.Now()
	orderID := fmt.Sprintf("qris-%s-%d", shortTxID, now.Unix())

	// Store order_id in payment entity for status checking
	paymentEntity.OrderID = orderID

	qrisReq := payment.QRISRequest{
		TransactionID:  transactionID,
		OrderID:        orderID,
		GrossAmount:    transaction.TotalAmount, // Use transaction total (includes tax & discount)
		CustomerName:   transaction.User.Name,
		CustomerEmail:  transaction.User.Email,
		Items:          uc.mapTransactionItemsToQRISItems(transaction),
		ExpiryDuration: uc.defaultExpiryMin,
	}

	qrisResponse, err := uc.midtransClient.GenerateQRIS(ctx, qrisReq)
	if err != nil {
		uc.logger.Error("Failed to generate new QRIS via Midtrans", "error", err)
		return nil, fmt.Errorf("failed to generate QRIS: %w", err)
	}

	// Update payment expiry using the same 'now' used for order_id
	newExpiry := now.Add(time.Duration(uc.defaultExpiryMin) * time.Minute)
	paymentEntity.ExpiresAt = newExpiry
	paymentEntity.Status = entities.PaymentPending
	paymentEntity.ExternalID = "" // Clear previous external ID
	paymentEntity.ExternalResponse = ""

	// Get existing QRIS code or create new one
	qrCodeEntity, err := uc.paymentRepo.GetQRISCodeByPaymentID(ctx, paymentEntity.ID)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
		// Create new QRIS code
		qrCodeEntity = entities.NewQRISCode(
			transactionID,
			paymentEntity.ID,
			qrisResponse.QRString,
			qrisResponse.URL,
			uc.defaultExpiryMin,
		)
	} else {
		// Update existing QRIS code
		qrCodeEntity.QRCode = qrisResponse.QRString
		qrCodeEntity.URL = qrisResponse.URL
		qrCodeEntity.ExpiresAt = newExpiry
	}

	// Update payment
	if err := uc.paymentRepo.UpdatePayment(ctx, paymentEntity); err != nil {
		uc.logger.Error("Failed to update payment", "error", err)
		return nil, err
	}

	// Update or create QRIS code
	if qrCodeEntity.ID == "" {
		// New QRIS code
		if err := uc.paymentRepo.CreateQRISCode(ctx, qrCodeEntity); err != nil {
			uc.logger.Error("Failed to create QRIS code", "error", err)
			return nil, err
		}
	} else {
		// Existing QRIS code
		if err := uc.paymentRepo.UpdateQRISCode(ctx, qrCodeEntity); err != nil {
			uc.logger.Error("Failed to update QRIS code", "error", err)
			return nil, err
		}
	}

	uc.logger.Info("QRIS refreshed successfully", "transaction_id", transactionID, "payment_id", paymentEntity.ID)

	return uc.mapPaymentToResponse(paymentEntity, qrCodeEntity), nil
}

// Helper methods
func (uc *PaymentUseCase) mapTransactionItemsToQRISItems(transaction *entities.Transaction) []payment.QRISItem {
	var qrisItems []payment.QRISItem

	// Add product items
	for _, item := range transaction.Items {
		qrisItems = append(qrisItems, payment.QRISItem{
			ID:       item.ProductID,
			Name:     item.Product.Name,
			Price:    item.UnitPrice,
			Quantity: item.Quantity,
		})
	}

	// Add tax as a line item if present
	if transaction.TaxAmount > 0 {
		qrisItems = append(qrisItems, payment.QRISItem{
			ID:       "TAX",
			Name:     "Tax",
			Price:    transaction.TaxAmount,
			Quantity: 1,
		})
	}

	// Add discount as negative line item if present
	if transaction.Discount > 0 {
		qrisItems = append(qrisItems, payment.QRISItem{
			ID:       "DISCOUNT",
			Name:     "Discount",
			Price:    -transaction.Discount, // Negative to reduce total
			Quantity: 1,
		})
	}

	return qrisItems
}

func (uc *PaymentUseCase) mapPaymentToResponse(payment *entities.Payment, qrisCode *entities.QRISCode) *PaymentResponse {
	response := &PaymentResponse{
		ID:            payment.ID,
		TransactionID: payment.TransactionID,
		Amount:        payment.Amount,
		Method:        payment.Method,
		Status:        payment.Status,
		ExternalID:    payment.ExternalID,
		ExpiresAt:     payment.ExpiresAt.Format(time.RFC3339),
		CreatedAt:     payment.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     payment.UpdatedAt.Format(time.RFC3339),
	}

	if payment.PaidAt != nil {
		paidAt := payment.PaidAt.Format(time.RFC3339)
		response.PaidAt = &paidAt
	}

	if qrisCode != nil {
		response.QRISCode = &QRISCodeResponse{
			ID:            qrisCode.ID,
			TransactionID: qrisCode.TransactionID,
			PaymentID:     qrisCode.PaymentID,
			QRCode:        qrisCode.QRCode,
			URL:           qrisCode.URL,
			ExpiresAt:     qrisCode.ExpiresAt.Format(time.RFC3339),
			CreatedAt:     qrisCode.CreatedAt.Format(time.RFC3339),
		}
	}

	return response
}
