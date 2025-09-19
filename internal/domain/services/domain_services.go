package services

import (
	"context"
	"errors"
	"goclean/internal/domain/entities"
	"goclean/internal/domain/repositories"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrProductNotFound    = errors.New("product not found")
	ErrOrderNotFound      = errors.New("order not found")
	ErrInvalidOrderStatus = errors.New("invalid order status")
)

// UserDomainService contains business logic for users
type UserDomainService struct {
	userRepo    repositories.UserRepository
	profileRepo repositories.ProfileRepository
}

// NewUserDomainService creates a new user domain service
func NewUserDomainService(userRepo repositories.UserRepository, profileRepo repositories.ProfileRepository) *UserDomainService {
	return &UserDomainService{
		userRepo:    userRepo,
		profileRepo: profileRepo,
	}
}

// CreateUserWithProfile creates a user and their profile
func (s *UserDomainService) CreateUserWithProfile(ctx context.Context, user *entities.User, profile *entities.Profile) error {
	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(ctx, user.Email)
	if existingUser != nil {
		return ErrUserAlreadyExists
	}

	// Create user
	if err := s.userRepo.Create(ctx, user); err != nil {
		return err
	}

	// Create profile if provided
	if profile != nil {
		profile.UserID = user.ID
		if err := s.profileRepo.Create(ctx, profile); err != nil {
			// Rollback user creation if profile creation fails
			s.userRepo.Delete(ctx, user.ID)
			return err
		}
	}

	return nil
}

// ProductDomainService contains business logic for products
type ProductDomainService struct {
	productRepo repositories.ProductRepository
}

// NewProductDomainService creates a new product domain service
func NewProductDomainService(productRepo repositories.ProductRepository) *ProductDomainService {
	return &ProductDomainService{
		productRepo: productRepo,
	}
}

// ValidateProduct validates product business rules
func (s *ProductDomainService) ValidateProduct(product *entities.Product) error {
	if product.Price <= 0 {
		return errors.New("product price must be greater than zero")
	}
	if product.Name == "" {
		return errors.New("product name is required")
	}
	if product.SKU == "" {
		return errors.New("product SKU is required")
	}
	return nil
}

// CreateProduct creates a product after validation
func (s *ProductDomainService) CreateProduct(ctx context.Context, product *entities.Product) error {
	if err := s.ValidateProduct(product); err != nil {
		return err
	}

	// Check if SKU already exists
	existingProduct, _ := s.productRepo.GetBySKU(ctx, product.SKU)
	if existingProduct != nil {
		return errors.New("product with this SKU already exists")
	}

	return s.productRepo.Create(ctx, product)
}

// OrderDomainService contains business logic for orders
type OrderDomainService struct {
	orderRepo   repositories.OrderRepository
	productRepo repositories.ProductRepository
}

// NewOrderDomainService creates a new order domain service
func NewOrderDomainService(orderRepo repositories.OrderRepository, productRepo repositories.ProductRepository) *OrderDomainService {
	return &OrderDomainService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

// CreateOrder creates an order with business validation
func (s *OrderDomainService) CreateOrder(ctx context.Context, order *entities.Order) error {
	if len(order.Items) == 0 {
		return errors.New("order must have at least one item")
	}

	totalPrice := 0.0
	for _, item := range order.Items {
		// Validate product exists
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return ErrProductNotFound
		}

		// Validate quantity
		if item.Quantity <= 0 {
			return errors.New("item quantity must be greater than zero")
		}

		// Set item price from product price
		item.Price = product.Price
		totalPrice += item.Price * float64(item.Quantity)
	}

	order.TotalPrice = totalPrice
	order.Status = entities.OrderStatusPending

	return s.orderRepo.Create(ctx, order)
}

// UpdateOrderStatus updates order status with business validation
func (s *OrderDomainService) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status entities.OrderStatus) error {
	if !status.IsValid() {
		return ErrInvalidOrderStatus
	}

	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return ErrOrderNotFound
	}

	// Business rules for status transitions
	if order.Status == entities.OrderStatusDelivered || order.Status == entities.OrderStatusCancelled {
		return errors.New("cannot change status of delivered or cancelled orders")
	}

	return s.orderRepo.UpdateStatus(ctx, orderID, status)
}
