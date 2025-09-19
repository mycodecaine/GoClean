package gorm

import (
	"context"
	"errors"
	"goclean/internal/domain/entities"
	"goclean/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// userRepository implements the UserRepository interface using GORM with soft delete support
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository with soft delete support
func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID gets a user by ID (excludes soft deleted)
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Preload("Profile").First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetByIDIncludeDeleted gets a user by ID (includes soft deleted)
func (r *userRepository) GetByIDIncludeDeleted(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Unscoped().Preload("Profile").First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail gets a user by email (excludes soft deleted)
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Preload("Profile").First(&user, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername gets a user by username (excludes soft deleted)
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Preload("Profile").First(&user, "username = ?", username).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete permanently deletes a user (hard delete)
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&entities.User{}, "id = ?", id).Error
}

// SoftDelete soft deletes a user
func (r *userRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.User{}, "id = ?", id).Error
}

// Restore restores a soft deleted user
func (r *userRepository) Restore(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entities.User{}).Unscoped().
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Update("deleted_at", nil).Error
}

// List lists users with pagination (excludes soft deleted)
func (r *userRepository) List(ctx context.Context, offset, limit int) ([]*entities.User, error) {
	var users []*entities.User
	err := r.db.WithContext(ctx).Preload("Profile").
		Offset(offset).Limit(limit).
		Find(&users).Error
	return users, err
}

// ListIncludeDeleted lists users with pagination (includes soft deleted)
func (r *userRepository) ListIncludeDeleted(ctx context.Context, offset, limit int) ([]*entities.User, error) {
	var users []*entities.User
	err := r.db.WithContext(ctx).Unscoped().Preload("Profile").
		Offset(offset).Limit(limit).
		Find(&users).Error
	return users, err
}

// ListDeleted lists only soft deleted users with pagination
func (r *userRepository) ListDeleted(ctx context.Context, offset, limit int) ([]*entities.User, error) {
	var users []*entities.User
	err := r.db.WithContext(ctx).Unscoped().Preload("Profile").
		Where("deleted_at IS NOT NULL").
		Offset(offset).Limit(limit).
		Find(&users).Error
	return users, err
}
