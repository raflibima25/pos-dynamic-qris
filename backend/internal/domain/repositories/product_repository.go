package repositories

import (
	"context"
	"qris-pos-backend/internal/domain/entities"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entities.Product) error
	GetByID(ctx context.Context, id string) (*entities.Product, error)
	GetBySKU(ctx context.Context, sku string) (*entities.Product, error)
	Update(ctx context.Context, product *entities.Product) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filters ProductFilters) ([]entities.Product, error)
	UpdateStock(ctx context.Context, id string, quantity int) error
	Search(ctx context.Context, query string, limit int) ([]entities.Product, error)
}

type ProductFilters struct {
	CategoryID string
	IsActive   *bool
	Limit      int
	Offset     int
}

type CategoryRepository interface {
	Create(ctx context.Context, category *entities.Category) error
	GetByID(ctx context.Context, id string) (*entities.Category, error)
	Update(ctx context.Context, category *entities.Category) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]entities.Category, error)
}