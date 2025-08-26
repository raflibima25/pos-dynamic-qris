package repositories

import (
	"context"
	"qris-pos-backend/internal/domain/entities"
	"qris-pos-backend/internal/domain/repositories"

	"gorm.io/gorm"
)

type productRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) repositories.ProductRepository {
	return &productRepositoryImpl{db: db}
}

func (r *productRepositoryImpl) Create(ctx context.Context, product *entities.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepositoryImpl) GetByID(ctx context.Context, id string) (*entities.Product, error) {
	var product entities.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Where("id = ?", id).
		First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepositoryImpl) GetBySKU(ctx context.Context, sku string) (*entities.Product, error) {
	var product entities.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Where("sku = ?", sku).
		First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepositoryImpl) Update(ctx context.Context, product *entities.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entities.Product{}, "id = ?", id).Error
}

func (r *productRepositoryImpl) List(ctx context.Context, filters repositories.ProductFilters) ([]entities.Product, error) {
	var products []entities.Product
	query := r.db.WithContext(ctx).Preload("Category")

	if filters.CategoryID != "" {
		query = query.Where("category_id = ?", filters.CategoryID)
	}

	if filters.IsActive != nil {
		query = query.Where("is_active = ?", *filters.IsActive)
	}

	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}

	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	err := query.Order("created_at DESC").Find(&products).Error
	return products, err
}

func (r *productRepositoryImpl) UpdateStock(ctx context.Context, id string, quantity int) error {
	return r.db.WithContext(ctx).
		Model(&entities.Product{}).
		Where("id = ?", id).
		Update("stock", gorm.Expr("stock + ?", quantity)).
		Error
}

func (r *productRepositoryImpl) Search(ctx context.Context, query string, limit int) ([]entities.Product, error) {
	var products []entities.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Where("name ILIKE ? OR sku ILIKE ?", "%"+query+"%", "%"+query+"%").
		Where("is_active = true").
		Limit(limit).
		Find(&products).Error
	return products, err
}

type categoryRepositoryImpl struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) repositories.CategoryRepository {
	return &categoryRepositoryImpl{db: db}
}

func (r *categoryRepositoryImpl) Create(ctx context.Context, category *entities.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *categoryRepositoryImpl) GetByID(ctx context.Context, id string) (*entities.Category, error) {
	var category entities.Category
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepositoryImpl) Update(ctx context.Context, category *entities.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *categoryRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entities.Category{}, "id = ?", id).Error
}

func (r *categoryRepositoryImpl) List(ctx context.Context, limit, offset int) ([]entities.Category, error) {
	var categories []entities.Category
	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("name ASC").
		Find(&categories).Error
	return categories, err
}