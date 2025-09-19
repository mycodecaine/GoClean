package persistence

import (
	"context"
	"goclean/internal/domain/entities"
	"goclean/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProfileGormRepository implements ProfileRepository using GORM
type ProfileGormRepository struct {
	db *gorm.DB
}

// NewProfileGormRepository creates a new profile GORM repository
func NewProfileGormRepository(db *gorm.DB) repositories.ProfileRepository {
	return &ProfileGormRepository{db: db}
}

// Create creates a new profile
func (r *ProfileGormRepository) Create(ctx context.Context, profile *entities.Profile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

// GetByID retrieves a profile by ID
func (r *ProfileGormRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Profile, error) {
	var profile entities.Profile
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// GetByIDIncludeDeleted gets a profile by ID including soft deleted
func (r *ProfileGormRepository) GetByIDIncludeDeleted(ctx context.Context, id uuid.UUID) (*entities.Profile, error) {
	var profile entities.Profile
	err := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// GetByUserID retrieves a profile by user ID
func (r *ProfileGormRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.Profile, error) {
	var profile entities.Profile
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// Update updates a profile
func (r *ProfileGormRepository) Update(ctx context.Context, profile *entities.Profile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}

// Delete deletes a profile (hard delete)
func (r *ProfileGormRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&entities.Profile{}, "id = ?", id).Error
}

// SoftDelete soft deletes a profile
func (r *ProfileGormRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Profile{}, "id = ?", id).Error
}

// Restore restores a soft deleted profile
func (r *ProfileGormRepository) Restore(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entities.Profile{}).Unscoped().
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Update("deleted_at", nil).Error
}
