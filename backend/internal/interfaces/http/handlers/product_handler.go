package handlers

import (
	"strconv"

	"qris-pos-backend/internal/usecases/product"
	"qris-pos-backend/pkg/logger"
	"qris-pos-backend/pkg/response"
	"qris-pos-backend/pkg/validator"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productUseCase *product.ProductUseCase
	logger         logger.Logger
}

func NewProductHandler(productUseCase *product.ProductUseCase, logger logger.Logger) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
		logger:         logger,
	}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product (Admin only)
// @Tags products
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body product.CreateProductRequest true "Product data"
// @Success 201 {object} response.Response{data=product.ProductResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req product.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format", err.Error())
		return
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.ValidationError(c, errors)
		return
	}

	result, err := h.productUseCase.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create product", "error", err)
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Product created successfully", result)
}

// GetProduct godoc
// @Summary Get product by ID
// @Description Get a single product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} response.Response{data=product.ProductResponse}
// @Failure 404 {object} response.Response
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")

	result, err := h.productUseCase.GetProduct(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get product", "error", err, "product_id", id)
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "Product retrieved successfully", result)
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update an existing product (Admin only)
// @Tags products
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Product ID"
// @Param request body product.UpdateProductRequest true "Updated product data"
// @Success 200 {object} response.Response{data=product.ProductResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var req product.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format", err.Error())
		return
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.ValidationError(c, errors)
		return
	}

	result, err := h.productUseCase.UpdateProduct(c.Request.Context(), id, &req)
	if err != nil {
		h.logger.Error("Failed to update product", "error", err, "product_id", id)
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Product updated successfully", result)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product (Admin only)
// @Tags products
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Product ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	err := h.productUseCase.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete product", "error", err, "product_id", id)
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Product deleted successfully", nil)
}

// ListProducts godoc
// @Summary List products
// @Description Get a list of products with optional filters
// @Tags products
// @Accept json
// @Produce json
// @Param category_id query string false "Filter by category ID"
// @Param is_active query boolean false "Filter by active status"
// @Param search query string false "Search in product name and SKU"
// @Param limit query int false "Number of products to return" default(20)
// @Param offset query int false "Number of products to skip" default(0)
// @Success 200 {object} response.Response{data=[]product.ProductResponse}
// @Router /products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	var filters product.ProductFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		response.BadRequest(c, "Invalid query parameters", err.Error())
		return
	}

	// Set defaults
	if filters.Limit == 0 {
		filters.Limit = 20
	}

	if errors := validator.ValidateStruct(filters); len(errors) > 0 {
		response.ValidationError(c, errors)
		return
	}

	result, err := h.productUseCase.ListProducts(c.Request.Context(), &filters)
	if err != nil {
		h.logger.Error("Failed to list products", "error", err)
		response.InternalError(c, "Failed to retrieve products", err.Error())
		return
	}

	response.Success(c, "Products retrieved successfully", result)
}

// UpdateStock godoc
// @Summary Update product stock
// @Description Update the stock quantity of a product
// @Tags products
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Product ID"
// @Param request body map[string]int true "Stock change" example({"quantity": 10})
// @Success 200 {object} response.Response{data=product.ProductResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /products/{id}/stock [patch]
func (h *ProductHandler) UpdateStock(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Quantity int `json:"quantity" validate:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format", err.Error())
		return
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.ValidationError(c, errors)
		return
	}

	result, err := h.productUseCase.UpdateStock(c.Request.Context(), id, req.Quantity)
	if err != nil {
		h.logger.Error("Failed to update product stock", "error", err, "product_id", id)
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Product stock updated successfully", result)
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new product category (Admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body product.CreateCategoryRequest true "Category data"
// @Success 201 {object} response.Response{data=product.CategoryResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /categories [post]
func (h *ProductHandler) CreateCategory(c *gin.Context) {
	var req product.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format", err.Error())
		return
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.ValidationError(c, errors)
		return
	}

	result, err := h.productUseCase.CreateCategory(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create category", "error", err)
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Category created successfully", result)
}

// ListCategories godoc
// @Summary List categories
// @Description Get a list of product categories
// @Tags categories
// @Accept json
// @Produce json
// @Param limit query int false "Number of categories to return" default(50)
// @Param offset query int false "Number of categories to skip" default(0)
// @Success 200 {object} response.Response{data=[]product.CategoryResponse}
// @Router /categories [get]
func (h *ProductHandler) ListCategories(c *gin.Context) {
	limit := 50
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	result, err := h.productUseCase.ListCategories(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Failed to list categories", "error", err)
		response.InternalError(c, "Failed to retrieve categories", err.Error())
		return
	}

	response.Success(c, "Categories retrieved successfully", result)
}