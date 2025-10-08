package handlers

import (
	"net/http"
	"qris-pos-backend/internal/usecases/payment"
	"qris-pos-backend/pkg/logger"
	"qris-pos-backend/pkg/response"
	"qris-pos-backend/pkg/validator"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentUseCase *payment.PaymentUseCase
	logger         logger.Logger
}

func NewPaymentHandler(paymentUseCase *payment.PaymentUseCase, logger logger.Logger) *PaymentHandler {
	return &PaymentHandler{
		paymentUseCase: paymentUseCase,
		logger:         logger,
	}
}

// GenerateQRIS godoc
// @Summary Generate QRIS for transaction
// @Description Generate a QRIS code for a pending transaction
// @Tags payments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body payment.GenerateQRISRequest true "QRIS generation data"
// @Success 201 {object} response.Response{data=payment.PaymentResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /payments/qris/generate [post]
func (h *PaymentHandler) GenerateQRIS(c *gin.Context) {
	var req payment.GenerateQRISRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format", err.Error())
		return
	}

	// Validate request
	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.ValidationError(c, errors)
		return
	}

	result, err := h.paymentUseCase.GenerateQRIS(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to generate QRIS", "error", err, "transaction_id", req.TransactionID)
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "QRIS generated successfully", result)
}

// GetPaymentStatus godoc
// @Summary Get payment status
// @Description Get the status of a payment for a transaction
// @Tags payments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param transaction_id path string true "Transaction ID"
// @Success 200 {object} response.Response{data=payment.PaymentStatusResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /payments/{transaction_id}/status [get]
func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	transactionID := c.Param("transaction_id")

	result, err := h.paymentUseCase.GetPaymentStatus(c.Request.Context(), transactionID)
	if err != nil {
		h.logger.Error("Failed to get payment status", "error", err, "transaction_id", transactionID)
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "Payment status retrieved successfully", result)
}

// RefreshQRIS godoc
// @Summary Refresh QRIS code
// @Description Refresh an expired QRIS code for a transaction
// @Tags payments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param transaction_id path string true "Transaction ID"
// @Success 200 {object} response.Response{data=payment.PaymentResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /qris/{transaction_id}/refresh [post]
func (h *PaymentHandler) RefreshQRIS(c *gin.Context) {
	transactionID := c.Param("transaction_id")
	if transactionID == "" {
		response.BadRequest(c, "Transaction ID is required", nil)
		return
	}

	result, err := h.paymentUseCase.RefreshQRIS(c.Request.Context(), transactionID)
	if err != nil {
		h.logger.Error("Failed to refresh QRIS", "error", err, "transaction_id", transactionID)
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "QRIS refreshed successfully", result)
}

// PaymentCallback godoc
// @Summary Payment callback from Midtrans
// @Description Handle payment notification from Midtrans
// @Tags payments
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "Midtrans notification data"
// @Success 200 {object} response.Response
// @Router /payments/callback [post]
func (h *PaymentHandler) PaymentCallback(c *gin.Context) {
	// Parse the notification data from Midtrans
	var notification map[string]interface{}
	if err := c.ShouldBindJSON(&notification); err != nil {
		h.logger.Error("Failed to parse payment callback", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Extract required fields
	orderID, ok := notification["order_id"].(string)
	if !ok {
		h.logger.Error("Missing order_id in payment callback")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing order_id"})
		return
	}

	status, ok := notification["transaction_status"].(string)
	if !ok {
		h.logger.Error("Missing transaction_status in payment callback")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing transaction_status"})
		return
	}

	externalID, _ := notification["transaction_id"].(string)
	responseData, _ := notification["response"].(string)

	// Handle the payment notification
	err := h.paymentUseCase.HandlePaymentNotification(c.Request.Context(), orderID, status, externalID, responseData)
	if err != nil {
		h.logger.Error("Failed to handle payment notification", "error", err, "order_id", orderID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process payment notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment notification processed successfully"})
}
