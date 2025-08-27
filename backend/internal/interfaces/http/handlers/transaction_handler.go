package handlers

import (
	"strconv"

	"qris-pos-backend/internal/domain/entities"
	"qris-pos-backend/internal/domain/repositories"
	"qris-pos-backend/internal/interfaces/middleware"
	"qris-pos-backend/internal/usecases/transaction"
	"qris-pos-backend/pkg/logger"
	"qris-pos-backend/pkg/response"
	"qris-pos-backend/pkg/validator"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionUseCase *transaction.TransactionUseCase
	logger             logger.Logger
}

func NewTransactionHandler(transactionUseCase *transaction.TransactionUseCase, logger logger.Logger) *TransactionHandler {
	return &TransactionHandler{
		transactionUseCase: transactionUseCase,
		logger:             logger,
	}
}

// CreateTransaction godoc
// @Summary Create a new transaction
// @Description Create a new transaction with items (shopping cart checkout)
// @Tags transactions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body transaction.CreateTransactionRequest true "Transaction data"
// @Success 201 {object} response.Response{data=transaction.TransactionResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req transaction.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format", err.Error())
		return
	}

	// Get current user from context
	currentUser, exists := middleware.GetCurrentUser(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Set user ID from authenticated user
	req.UserID = currentUser.UserID

	// Validate request
	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.ValidationError(c, errors)
		return
	}

	result, err := h.transactionUseCase.CreateTransaction(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create transaction", "error", err, "user_id", req.UserID)
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Transaction created successfully", result)
}

// GetTransaction godoc
// @Summary Get transaction by ID
// @Description Get a single transaction by its ID
// @Tags transactions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Transaction ID"
// @Success 200 {object} response.Response{data=transaction.TransactionResponse}
// @Failure 404 {object} response.Response
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	id := c.Param("id")

	result, err := h.transactionUseCase.GetTransaction(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get transaction", "error", err, "transaction_id", id)
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "Transaction retrieved successfully", result)
}

// ListTransactions godoc
// @Summary List transactions
// @Description Get a list of transactions with optional filters
// @Tags transactions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user_id query string false "Filter by user ID"
// @Param status query string false "Filter by status"
// @Param date_from query string false "Filter by date from (YYYY-MM-DD)"
// @Param date_to query string false "Filter by date to (YYYY-MM-DD)"
// @Param limit query int false "Number of transactions to return" default(20)
// @Param offset query int false "Number of transactions to skip" default(0)
// @Success 200 {object} response.Response{data=[]transaction.TransactionResponse}
// @Router /transactions [get]
func (h *TransactionHandler) ListTransactions(c *gin.Context) {
	filters := repositories.TransactionFilters{
		UserID: c.Query("user_id"),
		Limit:  20, // default
		Offset: 0,  // default
	}

	// Convert status string to enum if provided
	if statusStr := c.Query("status"); statusStr != "" {
		filters.Status = entities.TransactionStatus(statusStr)
	}

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		filters.DateFrom = &dateFrom
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		filters.DateTo = &dateTo
	}

	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			filters.Limit = l
		}
	}

	if offset := c.Query("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
			filters.Offset = o
		}
	}

	result, err := h.transactionUseCase.ListTransactions(c.Request.Context(), filters)
	if err != nil {
		h.logger.Error("Failed to list transactions", "error", err)
		response.InternalError(c, "Failed to retrieve transactions", err.Error())
		return
	}

	response.Success(c, "Transactions retrieved successfully", result)
}

// AddItemToTransaction godoc
// @Summary Add item to transaction
// @Description Add a product item to an existing pending transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Transaction ID"
// @Param request body transaction.AddItemRequest true "Item data"
// @Success 200 {object} response.Response{data=transaction.TransactionResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /transactions/{id}/items [post]
func (h *TransactionHandler) AddItemToTransaction(c *gin.Context) {
	id := c.Param("id")

	var req transaction.AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format", err.Error())
		return
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.ValidationError(c, errors)
		return
	}

	result, err := h.transactionUseCase.AddItemToTransaction(c.Request.Context(), id, &req)
	if err != nil {
		h.logger.Error("Failed to add item to transaction", "error", err, "transaction_id", id)
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Item added to transaction successfully", result)
}

// RemoveItemFromTransaction godoc
// @Summary Remove item from transaction
// @Description Remove a product item from an existing pending transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Transaction ID"
// @Param product_id path string true "Product ID"
// @Success 200 {object} response.Response{data=transaction.TransactionResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /transactions/{id}/items/{product_id} [delete]
func (h *TransactionHandler) RemoveItemFromTransaction(c *gin.Context) {
	id := c.Param("id")
	productID := c.Param("product_id")

	result, err := h.transactionUseCase.RemoveItemFromTransaction(c.Request.Context(), id, productID)
	if err != nil {
		h.logger.Error("Failed to remove item from transaction", "error", err, "transaction_id", id, "product_id", productID)
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Item removed from transaction successfully", result)
}

// UpdateItemQuantity godoc
// @Summary Update item quantity in transaction
// @Description Update the quantity of a product item in an existing pending transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Transaction ID"
// @Param product_id path string true "Product ID"
// @Param request body transaction.UpdateItemRequest true "Quantity data"
// @Success 200 {object} response.Response{data=transaction.TransactionResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /transactions/{id}/items/{product_id} [patch]
func (h *TransactionHandler) UpdateItemQuantity(c *gin.Context) {
	id := c.Param("id")
	productID := c.Param("product_id")

	var req transaction.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format", err.Error())
		return
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.ValidationError(c, errors)
		return
	}

	result, err := h.transactionUseCase.UpdateItemQuantity(c.Request.Context(), id, productID, &req)
	if err != nil {
		h.logger.Error("Failed to update item quantity", "error", err, "transaction_id", id, "product_id", productID)
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Item quantity updated successfully", result)
}

// CancelTransaction godoc
// @Summary Cancel a transaction
// @Description Cancel a pending transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Transaction ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /transactions/{id}/cancel [put]
func (h *TransactionHandler) CancelTransaction(c *gin.Context) {
	id := c.Param("id")

	err := h.transactionUseCase.CancelTransaction(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to cancel transaction", "error", err, "transaction_id", id)
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Transaction cancelled successfully", nil)
}
