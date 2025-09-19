package services

import (
	"context"
	"errors"
	"goclean/internal/domain/entities"
	"goclean/internal/domain/events"
	"goclean/internal/domain/repositories"
	"goclean/pkg/logger"

	"github.com/google/uuid"
)

// UserAggregateService handles User aggregate operations with domain events and soft delete
type UserAggregateService struct {
	userRepo        repositories.UserRepository
	profileRepo     repositories.ProfileRepository
	eventDispatcher *events.DomainEventDispatcher
	logger          *logger.Logger
}

// NewUserAggregateService creates a new user aggregate service
func NewUserAggregateService(
	userRepo repositories.UserRepository,
	profileRepo repositories.ProfileRepository,
	eventDispatcher *events.DomainEventDispatcher,
	logger *logger.Logger,
) *UserAggregateService {
	return &UserAggregateService{
		userRepo:        userRepo,
		profileRepo:     profileRepo,
		eventDispatcher: eventDispatcher,
		logger:          logger,
	}
}

// CreateUser creates a new user using aggregate root pattern
func (s *UserAggregateService) CreateUser(ctx context.Context, email, username, firstName, lastName string) (*entities.User, error) {
	// Create user aggregate using factory method
	user := entities.NewUser(email, username, firstName, lastName)

	// Validate business rules
	if err := s.validateUserBusinessRules(ctx, user); err != nil {
		return nil, err
	}

	// Persist the aggregate
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Dispatch domain events
	if err := s.eventDispatcher.DispatchEvents(ctx, &user.AggregateRoot); err != nil {
		s.logger.Error("Failed to dispatch domain events", "error", err)
		// Note: In production, you might want to handle this more gracefully
		// Perhaps by storing events for later processing
	}

	return user, nil
}

// SoftDeleteUser performs soft delete on user aggregate
func (s *UserAggregateService) SoftDeleteUser(ctx context.Context, userID uuid.UUID) error {
	// Get user aggregate
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Validate business rules for deletion
	if err := s.validateUserDeletionRules(ctx, user); err != nil {
		return err
	}

	// Perform soft delete through aggregate method
	user.Delete()

	// Update the aggregate
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Dispatch domain events
	if err := s.eventDispatcher.DispatchEvents(ctx, &user.AggregateRoot); err != nil {
		s.logger.Error("Failed to dispatch domain events", "error", err)
	}

	return nil
}

// RestoreUser restores a soft deleted user
func (s *UserAggregateService) RestoreUser(ctx context.Context, userID uuid.UUID) error {
	// Get user including deleted ones
	user, err := s.userRepo.GetByIDIncludeDeleted(ctx, userID)
	if err != nil {
		return err
	}

	if !user.IsDeleted() {
		return errors.New("user is not deleted")
	}

	// Restore through aggregate method
	user.Restore()
	user.IsActive = true

	// Update the aggregate
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

// GetUserWithDeleted gets a user including soft deleted ones
func (s *UserAggregateService) GetUserWithDeleted(ctx context.Context, userID uuid.UUID) (*entities.User, error) {
	return s.userRepo.GetByIDIncludeDeleted(ctx, userID)
}

// ListDeletedUsers lists all soft deleted users
func (s *UserAggregateService) ListDeletedUsers(ctx context.Context, offset, limit int) ([]*entities.User, error) {
	return s.userRepo.ListDeleted(ctx, offset, limit)
}

// validateUserBusinessRules validates business rules for user creation
func (s *UserAggregateService) validateUserBusinessRules(ctx context.Context, user *entities.User) error {
	// Check if email already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return errors.New("email already exists")
	}

	// Check if username already exists
	existingUser, err = s.userRepo.GetByUsername(ctx, user.Username)
	if err == nil && existingUser != nil {
		return errors.New("username already exists")
	}

	// Additional business rules can be added here
	return nil
}

// validateUserDeletionRules validates business rules for user deletion
func (s *UserAggregateService) validateUserDeletionRules(ctx context.Context, user *entities.User) error {
	if user.IsDeleted() {
		return errors.New("user is already deleted")
	}

	// You can add more business rules here, for example:
	// - Check if user has active subscriptions
	// - Check if user has pending orders
	// - Check if user is an admin

	return nil
}
