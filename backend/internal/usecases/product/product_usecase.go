package product

import (
	"context"
	"errors"

	"qris-pos-backend/internal/domain/entities"
	"qris-pos-backend/internal/domain/repositories"
	appErrors "qris-pos-backend/pkg/errors"
	"qris-pos-backend/pkg/logger"

	"gorm.io/gorm"
)

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gte=0"`
	Stock       int     `json:"stock" validate:"required,gte=0"`
	CategoryID  string  `json:"category_id" validate:"required,uuid"`
	SKU         string  `json:"sku"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gte=0"`
	Stock       int     `json:"stock" validate:"required,gte=0"`
	CategoryID  string  `json:"category_id" validate:"required,uuid"`
	SKU         string  `json:"sku"`
	IsActive    *bool   `json:"is_active"`
}

type ProductResponse struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Price       float64                `json:"price"`
	Stock       int                    `json:"stock"`
	CategoryID  string                 `json:"category_id"`
	SKU         string                 `json:"sku"`
	IsActive    bool                   `json:"is_active"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
	Category    *CategoryResponse      `json:"category,omitempty"`
}

type CategoryResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

type UpdateCategoryRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=255"`
	IsActive *bool  `json:"is_active"`
}

type ProductFilters struct {
	CategoryID string `form:"category_id"`
	IsActive   *bool  `form:"is_active"`
	Search     string `form:"search"`
	Limit      int    `form:"limit,default=20" validate:"gte=1,lte=100"`
	Offset     int    `form:"offset,default=0" validate:"gte=0"`
}

type ProductUseCase struct {
	productRepo  repositories.ProductRepository
	categoryRepo repositories.CategoryRepository
	logger       logger.Logger
}

func NewProductUseCase(
	productRepo repositories.ProductRepository,
	categoryRepo repositories.CategoryRepository,
	logger logger.Logger,
) *ProductUseCase {
	return &ProductUseCase{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (uc *ProductUseCase) CreateProduct(ctx context.Context, req *CreateProductRequest) (*ProductResponse, error) {
	// Validate category exists
	_, err := uc.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	// Check if SKU already exists (if provided)
	if req.SKU != "" {
		existingProduct, err := uc.productRepo.GetBySKU(ctx, req.SKU)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existingProduct != nil {
			return nil, appErrors.ErrSKUExists
		}
	}

	product, err := entities.NewProduct(req.Name, req.Description, req.SKU, req.CategoryID, req.Price, req.Stock)
	if err != nil {
		return nil, err
	}

	if err := uc.productRepo.Create(ctx, product); err != nil {
		uc.logger.Error("Failed to create product", "error", err)
		return nil, err
	}

	// Get product with category
	createdProduct, err := uc.productRepo.GetByID(ctx, product.ID)
	if err != nil {
		return nil, err
	}

	uc.logger.Info("Product created successfully", "product_id", product.ID, "name", product.Name)
	return uc.mapProductToResponse(createdProduct), nil
}

func (uc *ProductUseCase) GetProduct(ctx context.Context, id string) (*ProductResponse, error) {
	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrProductNotFound
		}
		return nil, err
	}

	return uc.mapProductToResponse(product), nil
}

func (uc *ProductUseCase) UpdateProduct(ctx context.Context, id string, req *UpdateProductRequest) (*ProductResponse, error) {
	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrProductNotFound
		}
		return nil, err
	}

	// Validate category exists
	_, err = uc.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	// Check if SKU already exists (if changed and provided)
	if req.SKU != "" && req.SKU != product.SKU {
		existingProduct, err := uc.productRepo.GetBySKU(ctx, req.SKU)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existingProduct != nil && existingProduct.ID != id {
			return nil, appErrors.ErrSKUExists
		}
	}

	// Update product fields
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock
	product.CategoryID = req.CategoryID
	product.SKU = req.SKU

	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := uc.productRepo.Update(ctx, product); err != nil {
		uc.logger.Error("Failed to update product", "error", err, "product_id", id)
		return nil, err
	}

	// Get updated product with category
	updatedProduct, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	uc.logger.Info("Product updated successfully", "product_id", id)
	return uc.mapProductToResponse(updatedProduct), nil
}

func (uc *ProductUseCase) DeleteProduct(ctx context.Context, id string) error {
	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return appErrors.ErrProductNotFound
		}
		return err
	}

	if err := uc.productRepo.Delete(ctx, id); err != nil {
		uc.logger.Error("Failed to delete product", "error", err, "product_id", id)
		return err
	}

	uc.logger.Info("Product deleted successfully", "product_id", id, "name", product.Name)
	return nil
}

func (uc *ProductUseCase) ListProducts(ctx context.Context, filters *ProductFilters) ([]ProductResponse, error) {
	repoFilters := repositories.ProductFilters{
		CategoryID: filters.CategoryID,
		IsActive:   filters.IsActive,
		Limit:      filters.Limit,
		Offset:     filters.Offset,
	}

	var products []entities.Product
	var err error

	if filters.Search != "" {
		products, err = uc.productRepo.Search(ctx, filters.Search, filters.Limit)
	} else {
		products, err = uc.productRepo.List(ctx, repoFilters)
	}

	if err != nil {
		uc.logger.Error("Failed to list products", "error", err)
		return nil, err
	}

	responses := make([]ProductResponse, len(products))
	for i, product := range products {
		responses[i] = *uc.mapProductToResponse(&product)
	}

	return responses, nil
}

func (uc *ProductUseCase) UpdateStock(ctx context.Context, id string, quantity int) (*ProductResponse, error) {
	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrProductNotFound
		}
		return nil, err
	}

	if err := product.UpdateStock(quantity); err != nil {
		return nil, err
	}

	if err := uc.productRepo.Update(ctx, product); err != nil {
		uc.logger.Error("Failed to update product stock", "error", err, "product_id", id)
		return nil, err
	}

	uc.logger.Info("Product stock updated", "product_id", id, "quantity_change", quantity, "new_stock", product.Stock)
	return uc.mapProductToResponse(product), nil
}

// Category operations
func (uc *ProductUseCase) CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*CategoryResponse, error) {
	category := &entities.Category{
		Name:     req.Name,
		IsActive: true,
	}

	if err := uc.categoryRepo.Create(ctx, category); err != nil {
		uc.logger.Error("Failed to create category", "error", err)
		return nil, err
	}

	uc.logger.Info("Category created successfully", "category_id", category.ID, "name", category.Name)
	return uc.mapCategoryToResponse(category), nil
}

func (uc *ProductUseCase) ListCategories(ctx context.Context, limit, offset int) ([]CategoryResponse, error) {
	categories, err := uc.categoryRepo.List(ctx, limit, offset)
	if err != nil {
		uc.logger.Error("Failed to list categories", "error", err)
		return nil, err
	}

	responses := make([]CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = *uc.mapCategoryToResponse(&category)
	}

	return responses, nil
}

func (uc *ProductUseCase) mapProductToResponse(product *entities.Product) *ProductResponse {
	response := &ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		CategoryID:  product.CategoryID,
		SKU:         product.SKU,
		IsActive:    product.IsActive,
		CreatedAt:   product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if product.Category.ID != "" {
		response.Category = uc.mapCategoryToResponse(&product.Category)
	}

	return response
}

func (uc *ProductUseCase) mapCategoryToResponse(category *entities.Category) *CategoryResponse {
	return &CategoryResponse{
		ID:       category.ID,
		Name:     category.Name,
		IsActive: category.IsActive,
	}
}