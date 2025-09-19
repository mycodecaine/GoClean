package test

import (
	"context"
	"goclean/internal/domain/entities"
	"goclean/internal/domain/services"
	"goclean/test/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserDomainService_CreateUserWithProfile(t *testing.T) {
	tests := []struct {
		name        string
		user        *entities.User
		profile     *entities.Profile
		setupMocks  func(*mocks.MockUserRepository, *mocks.MockProfileRepository)
		expectError bool
		errorMsg    string
	}{
		{
			name:    "successful user creation with profile",
			user:    entities.NewUser("test@example.com", "testuser", "Test", "User"),
			profile: entities.NewProfile(uuid.New(), "Test bio", "avatar.jpg", nil),
			setupMocks: func(userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)
				userRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.User")).Return(nil)
				profileRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.Profile")).Return(nil)
			},
			expectError: false,
		},
		{
			name:    "user already exists",
			user:    entities.NewUser("existing@example.com", "existinguser", "Existing", "User"),
			profile: nil,
			setupMocks: func(userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				existingUser := entities.NewUser("existing@example.com", "existinguser", "Existing", "User")
				userRepo.On("GetByEmail", mock.Anything, "existing@example.com").Return(existingUser, nil)
			},
			expectError: true,
			errorMsg:    "user already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			userRepo := &mocks.MockUserRepository{}
			profileRepo := &mocks.MockProfileRepository{}

			if tt.setupMocks != nil {
				tt.setupMocks(userRepo, profileRepo)
			}

			// Create service
			service := services.NewUserDomainService(userRepo, profileRepo)

			// Execute
			err := service.CreateUserWithProfile(context.Background(), tt.user, tt.profile)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			// Verify mocks
			userRepo.AssertExpectations(t)
			profileRepo.AssertExpectations(t)
		})
	}
}

func TestProductDomainService_ValidateProduct(t *testing.T) {
	tests := []struct {
		name        string
		product     *entities.Product
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid product",
			product: &entities.Product{
				Name:  "Test Product",
				Price: 19.99,
				SKU:   "TEST-001",
			},
			expectError: false,
		},
		{
			name: "invalid price - zero",
			product: &entities.Product{
				Name:  "Test Product",
				Price: 0,
				SKU:   "TEST-001",
			},
			expectError: true,
			errorMsg:    "price must be greater than zero",
		},
		{
			name: "invalid price - negative",
			product: &entities.Product{
				Name:  "Test Product",
				Price: -10.00,
				SKU:   "TEST-001",
			},
			expectError: true,
			errorMsg:    "price must be greater than zero",
		},
		{
			name: "missing name",
			product: &entities.Product{
				Name:  "",
				Price: 19.99,
				SKU:   "TEST-001",
			},
			expectError: true,
			errorMsg:    "name is required",
		},
		{
			name: "missing SKU",
			product: &entities.Product{
				Name:  "Test Product",
				Price: 19.99,
				SKU:   "",
			},
			expectError: true,
			errorMsg:    "SKU is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			productRepo := &mocks.MockProductRepository{}
			service := services.NewProductDomainService(productRepo)

			// Execute
			err := service.ValidateProduct(tt.product)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOrderStatus_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		status   entities.OrderStatus
		expected bool
	}{
		{"pending status", entities.OrderStatusPending, true},
		{"confirmed status", entities.OrderStatusConfirmed, true},
		{"shipped status", entities.OrderStatusShipped, true},
		{"delivered status", entities.OrderStatusDelivered, true},
		{"cancelled status", entities.OrderStatusCancelled, true},
		{"invalid status", entities.OrderStatus("invalid"), false},
		{"empty status", entities.OrderStatus(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}
