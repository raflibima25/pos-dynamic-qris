package payment

import (
	"context"
	"fmt"
	"qris-pos-backend/internal/infrastructure/config"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

// MidtransClient wraps the Midtrans SDK client
type MidtransClient struct {
	coreAPIClient *coreapi.Client
	config        config.MidtransConfig
}

// NewMidtransClient creates a new Midtrans client instance
func NewMidtransClient(cfg config.MidtransConfig) *MidtransClient {
	coreAPIClient := &coreapi.Client{}
	coreAPIClient.New(cfg.ServerKey, getEnvironment(cfg.Environment))

	return &MidtransClient{
		coreAPIClient: coreAPIClient,
		config:        cfg,
	}
}

// Helper function to get environment
func getEnvironment(env string) midtrans.EnvironmentType {
	if env == "production" {
		return midtrans.Production
	}
	return midtrans.Sandbox
}

// QRISRequest represents the data needed to generate a QRIS code
type QRISRequest struct {
	TransactionID   string
	OrderID         string
	GrossAmount     float64
	CustomerName    string
	CustomerEmail   string
	CustomerPhone   string
	Items           []QRISItem
	ExpiryDuration  int // in minutes
}

// QRISItem represents an item in the QRIS transaction
type QRISItem struct {
	ID       string
	Name     string
	Price    float64
	Quantity int
}

// QRISResponse represents the response from Midtrans
type QRISResponse struct {
	Token    string
	QRString string
	URL      string // Simulator URL for testing
}

// GenerateQRIS generates a QRIS code for payment
func (m *MidtransClient) GenerateQRIS(ctx context.Context, req QRISRequest) (*QRISResponse, error) {
	// Check context cancellation
	if ctx.Err() != nil {
		return nil, fmt.Errorf("context cancelled: %w", ctx.Err())
	}

	// Prepare items
	var items []midtrans.ItemDetails
	for _, item := range req.Items {
		items = append(items, midtrans.ItemDetails{
			ID:    item.ID,
			Name:  item.Name,
			Price: int64(item.Price), // Price already in correct format (IDR)
			Qty:   int32(item.Quantity),
		})
	}

	// Create charge request for QRIS using map approach
	chargeReq := &coreapi.ChargeReqWithMap{
		"payment_type": "qris",
		"transaction_details": map[string]interface{}{
			"order_id":     req.OrderID,
			"gross_amount": int64(req.GrossAmount), // Amount already in correct format (IDR)
		},
		"item_details": items,
		"customer_details": map[string]interface{}{
			"first_name": req.CustomerName,
			"email":      req.CustomerEmail,
			"phone":      req.CustomerPhone,
		},
	}

	// Charge the transaction
	res, err := m.coreAPIClient.ChargeTransactionWithMap(chargeReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create Midtrans transaction: %w", err)
	}

	// Extract QR string from response
	// Use qr_string field which contains the actual QRIS string for QR code generation
	var qrString string
	if qrStr, ok := res["qr_string"].(string); ok {
		qrString = qrStr
	}

	// Extract simulator URL from actions
	var simulatorURL string
	if actions, ok := res["actions"].([]interface{}); ok && len(actions) > 0 {
		if action, ok := actions[0].(map[string]interface{}); ok {
			if url, ok := action["url"].(string); ok {
				simulatorURL = url
			}
		}
	}

	// Extract transaction ID
	token := ""
	if transactionID, ok := res["transaction_id"].(string); ok {
		token = transactionID
	}

	return &QRISResponse{
		Token:    token,
		QRString: qrString,
		URL:      simulatorURL,
	}, nil
}

// GetTransactionStatus gets the status of a transaction
func (m *MidtransClient) GetTransactionStatus(ctx context.Context, orderID string) (*coreapi.TransactionStatusResponse, error) {
	res, err := m.coreAPIClient.CheckTransaction(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to check transaction status: %w", err)
	}
	return res, nil
}

// CancelTransaction cancels a transaction
func (m *MidtransClient) CancelTransaction(ctx context.Context, orderID string) error {
	_, err := m.coreAPIClient.CancelTransaction(orderID)
	if err != nil {
		return fmt.Errorf("failed to cancel transaction: %w", err)
	}
	return nil
}
