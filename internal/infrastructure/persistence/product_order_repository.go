package persistence

import (
	"context"
	"goclean/internal/domain/entities"
	"goclean/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductGormRepository implements ProductRepository using GORM
type ProductGormRepository struct {
	db *gorm.DB
}

// NewProductGormRepository creates a new product GORM repository
func NewProductGormRepository(db *gorm.DB) repositories.ProductRepository {
	return &ProductGormRepository{db: db}
}

// Create creates a new product
func (r *ProductGormRepository) Create(ctx context.Context, product *entities.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// GetByID retrieves a product by ID
func (r *ProductGormRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	var product entities.Product
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetBySKU retrieves a product by SKU
func (r *ProductGormRepository) GetBySKU(ctx context.Context, sku string) (*entities.Product, error) {
	var product entities.Product
	err := r.db.WithContext(ctx).Where("sku = ?", sku).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Update updates a product
func (r *ProductGormRepository) Update(ctx context.Context, product *entities.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

// Delete deletes a product (hard delete)
func (r *ProductGormRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&entities.Product{}, "id = ?", id).Error
}

// GetByIDIncludeDeleted gets a product by ID including soft deleted
func (r *ProductGormRepository) GetByIDIncludeDeleted(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	var product entities.Product
	err := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// SoftDelete soft deletes a product
func (r *ProductGormRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Product{}, "id = ?", id).Error
}

// Restore restores a soft deleted product
func (r *ProductGormRepository) Restore(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entities.Product{}).Unscoped().
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Update("deleted_at", nil).Error
}

// ListIncludeDeleted lists products including soft deleted
func (r *ProductGormRepository) ListIncludeDeleted(ctx context.Context, offset, limit int) ([]*entities.Product, error) {
	var products []*entities.Product
	err := r.db.WithContext(ctx).Unscoped().Offset(offset).Limit(limit).Find(&products).Error
	return products, err
}

// ListDeleted lists only soft deleted products
func (r *ProductGormRepository) ListDeleted(ctx context.Context, offset, limit int) ([]*entities.Product, error) {
	var products []*entities.Product
	err := r.db.WithContext(ctx).Unscoped().
		Where("deleted_at IS NOT NULL").
		Offset(offset).Limit(limit).Find(&products).Error
	return products, err
}

// List retrieves products with pagination
func (r *ProductGormRepository) List(ctx context.Context, offset, limit int) ([]*entities.Product, error) {
	var products []*entities.Product
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&products).Error
	return products, err
}

// ListByCategory retrieves products by category with pagination
func (r *ProductGormRepository) ListByCategory(ctx context.Context, category string, offset, limit int) ([]*entities.Product, error) {
	var products []*entities.Product
	err := r.db.WithContext(ctx).Where("category = ?", category).Offset(offset).Limit(limit).Find(&products).Error
	return products, err
}

// Search searches products by name or description with pagination
func (r *ProductGormRepository) Search(ctx context.Context, query string, offset, limit int) ([]*entities.Product, error) {
	var products []*entities.Product
	searchQuery := "%" + query + "%"
	err := r.db.WithContext(ctx).Where("name ILIKE ? OR description ILIKE ?", searchQuery, searchQuery).
		Offset(offset).Limit(limit).Find(&products).Error
	return products, err
}

// OrderGormRepository implements OrderRepository using GORM
type OrderGormRepository struct {
	db *gorm.DB
}

// NewOrderGormRepository creates a new order GORM repository
func NewOrderGormRepository(db *gorm.DB) repositories.OrderRepository {
	return &OrderGormRepository{db: db}
}

// Create creates a new order
func (r *OrderGormRepository) Create(ctx context.Context, order *entities.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// GetByID retrieves an order by ID
func (r *OrderGormRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Order, error) {
	var order entities.Order
	err := r.db.WithContext(ctx).Preload("Items").Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByUserID retrieves orders by user ID with pagination
func (r *OrderGormRepository) GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*entities.Order, error) {
	var orders []*entities.Order
	err := r.db.WithContext(ctx).Preload("Items").Where("user_id = ?", userID).
		Offset(offset).Limit(limit).Find(&orders).Error
	return orders, err
}

// Update updates an order
func (r *OrderGormRepository) Update(ctx context.Context, order *entities.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

// Delete deletes an order (hard delete)
func (r *OrderGormRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&entities.Order{}, "id = ?", id).Error
}

// GetByIDIncludeDeleted gets an order by ID including soft deleted
func (r *OrderGormRepository) GetByIDIncludeDeleted(ctx context.Context, id uuid.UUID) (*entities.Order, error) {
	var order entities.Order
	err := r.db.WithContext(ctx).Unscoped().Preload("Items").Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// SoftDelete soft deletes an order
func (r *OrderGormRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Order{}, "id = ?", id).Error
}

// Restore restores a soft deleted order
func (r *OrderGormRepository) Restore(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entities.Order{}).Unscoped().
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Update("deleted_at", nil).Error
}

// ListIncludeDeleted lists orders including soft deleted
func (r *OrderGormRepository) ListIncludeDeleted(ctx context.Context, offset, limit int) ([]*entities.Order, error) {
	var orders []*entities.Order
	err := r.db.WithContext(ctx).Unscoped().Preload("Items").Offset(offset).Limit(limit).Find(&orders).Error
	return orders, err
}

// ListDeleted lists only soft deleted orders
func (r *OrderGormRepository) ListDeleted(ctx context.Context, offset, limit int) ([]*entities.Order, error) {
	var orders []*entities.Order
	err := r.db.WithContext(ctx).Unscoped().Preload("Items").
		Where("deleted_at IS NOT NULL").
		Offset(offset).Limit(limit).Find(&orders).Error
	return orders, err
}

// List retrieves orders with pagination
func (r *OrderGormRepository) List(ctx context.Context, offset, limit int) ([]*entities.Order, error) {
	var orders []*entities.Order
	err := r.db.WithContext(ctx).Preload("Items").Offset(offset).Limit(limit).Find(&orders).Error
	return orders, err
}

// UpdateStatus updates order status
func (r *OrderGormRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entities.OrderStatus) error {
	return r.db.WithContext(ctx).Model(&entities.Order{}).Where("id = ?", id).Update("status", status).Error
}
