package repositories

import (
	"context"
	"goclean/internal/domain/entities"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByIDIncludeDeleted(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id uuid.UUID) error     // Hard delete
	SoftDelete(ctx context.Context, id uuid.UUID) error // Soft delete
	Restore(ctx context.Context, id uuid.UUID) error    // Restore soft deleted
	List(ctx context.Context, offset, limit int) ([]*entities.User, error)
	ListIncludeDeleted(ctx context.Context, offset, limit int) ([]*entities.User, error)
	ListDeleted(ctx context.Context, offset, limit int) ([]*entities.User, error)
}

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	Create(ctx context.Context, product *entities.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Product, error)
	GetByIDIncludeDeleted(ctx context.Context, id uuid.UUID) (*entities.Product, error)
	GetBySKU(ctx context.Context, sku string) (*entities.Product, error)
	Update(ctx context.Context, product *entities.Product) error
	Delete(ctx context.Context, id uuid.UUID) error     // Hard delete
	SoftDelete(ctx context.Context, id uuid.UUID) error // Soft delete
	Restore(ctx context.Context, id uuid.UUID) error    // Restore soft deleted
	List(ctx context.Context, offset, limit int) ([]*entities.Product, error)
	ListIncludeDeleted(ctx context.Context, offset, limit int) ([]*entities.Product, error)
	ListDeleted(ctx context.Context, offset, limit int) ([]*entities.Product, error)
	ListByCategory(ctx context.Context, category string, offset, limit int) ([]*entities.Product, error)
	Search(ctx context.Context, query string, offset, limit int) ([]*entities.Product, error)
}

// OrderRepository defines the interface for order data access
type OrderRepository interface {
	Create(ctx context.Context, order *entities.Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Order, error)
	GetByIDIncludeDeleted(ctx context.Context, id uuid.UUID) (*entities.Order, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*entities.Order, error)
	Update(ctx context.Context, order *entities.Order) error
	Delete(ctx context.Context, id uuid.UUID) error     // Hard delete
	SoftDelete(ctx context.Context, id uuid.UUID) error // Soft delete
	Restore(ctx context.Context, id uuid.UUID) error    // Restore soft deleted
	List(ctx context.Context, offset, limit int) ([]*entities.Order, error)
	ListIncludeDeleted(ctx context.Context, offset, limit int) ([]*entities.Order, error)
	ListDeleted(ctx context.Context, offset, limit int) ([]*entities.Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status entities.OrderStatus) error
}

// ProfileRepository defines the interface for profile data access
type ProfileRepository interface {
	Create(ctx context.Context, profile *entities.Profile) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Profile, error)
	GetByIDIncludeDeleted(ctx context.Context, id uuid.UUID) (*entities.Profile, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.Profile, error)
	Update(ctx context.Context, profile *entities.Profile) error
	Delete(ctx context.Context, id uuid.UUID) error     // Hard delete
	SoftDelete(ctx context.Context, id uuid.UUID) error // Soft delete
	Restore(ctx context.Context, id uuid.UUID) error    // Restore soft deleted
}
