package dto

import (
	"time"

	"github.com/google/uuid"
)

// UserDTO represents user data transfer object
type UserDTO struct {
	ID        uuid.UUID   `json:"id"`
	Email     string      `json:"email"`
	Username  string      `json:"username"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	IsActive  bool        `json:"is_active"`
	Profile   *ProfileDTO `json:"profile,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// ProfileDTO represents profile data transfer object
type ProfileDTO struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	Bio         string     `json:"bio"`
	Avatar      string     `json:"avatar"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ProductDTO represents product data transfer object
type ProductDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	SKU         string    `json:"sku"`
	Category    string    `json:"category"`
	IsActive    bool      `json:"is_active"`
	CreatedBy   uuid.UUID `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// OrderDTO represents order data transfer object
type OrderDTO struct {
	ID         uuid.UUID      `json:"id"`
	UserID     uuid.UUID      `json:"user_id"`
	Status     string         `json:"status"`
	TotalPrice float64        `json:"total_price"`
	Items      []OrderItemDTO `json:"items"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

// OrderItemDTO represents order item data transfer object
type OrderItemDTO struct {
	ID        uuid.UUID `json:"id"`
	OrderID   uuid.UUID `json:"order_id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateUserRequest represents create user request
type CreateUserRequest struct {
	Email     string                `json:"email" validate:"required,email"`
	Username  string                `json:"username" validate:"required,min=3,max=50"`
	FirstName string                `json:"first_name" validate:"required,min=1,max=100"`
	LastName  string                `json:"last_name" validate:"required,min=1,max=100"`
	Profile   *CreateProfileRequest `json:"profile,omitempty"`
}

// CreateProfileRequest represents create profile request
type CreateProfileRequest struct {
	Bio         string `json:"bio"`
	Avatar      string `json:"avatar"`
	DateOfBirth string `json:"date_of_birth"`
}

// UpdateUserRequest represents update user request
type UpdateUserRequest struct {
	Email     *string `json:"email,omitempty" validate:"omitempty,email"`
	Username  *string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=1,max=100"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,min=1,max=100"`
}

// CreateProductRequest represents create product request
type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	SKU         string  `json:"sku" validate:"required,min=1,max=100"`
	Category    string  `json:"category" validate:"required"`
}

// UpdateProductRequest represents update product request
type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,gt=0"`
	Category    *string  `json:"category,omitempty"`
	IsActive    *bool    `json:"is_active,omitempty"`
}

// CreateOrderRequest represents create order request
type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items" validate:"required,min=1"`
}

// CreateOrderItemRequest represents create order item request
type CreateOrderItemRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,gt=0"`
}

// UpdateOrderStatusRequest represents update order status request
type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=pending confirmed shipped delivered cancelled"`
}

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Offset int `query:"offset" validate:"min=0"`
	Limit  int `query:"limit" validate:"min=1,max=100"`
}

// APIResponse represents standard API response
type APIResponse[T any] struct {
	Success bool   `json:"success"`
	Data    T      `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

// PaginatedResponse represents paginated API response
type PaginatedResponse[T any] struct {
	APIResponse[T]
	Pagination PaginationInfo `json:"pagination"`
}

// PaginationInfo represents pagination information
type PaginationInfo struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
}

// Swagger-compatible response types (non-generic versions for documentation)

// UserAPIResponse represents API response for user operations
type UserAPIResponse struct {
	Success bool     `json:"success"`
	Data    *UserDTO `json:"data,omitempty"`
	Error   string   `json:"error,omitempty"`
	Message string   `json:"message,omitempty"`
}

// ProductAPIResponse represents API response for product operations
type ProductAPIResponse struct {
	Success bool        `json:"success"`
	Data    *ProductDTO `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// OrderAPIResponse represents API response for order operations
type OrderAPIResponse struct {
	Success bool      `json:"success"`
	Data    *OrderDTO `json:"data,omitempty"`
	Error   string    `json:"error,omitempty"`
	Message string    `json:"message,omitempty"`
}

// ErrorAPIResponse represents API response for errors
type ErrorAPIResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

// UsersListResponse represents API response for user list operations
type UsersListResponse struct {
	Success bool      `json:"success"`
	Data    []UserDTO `json:"data,omitempty"`
	Error   string    `json:"error,omitempty"`
	Message string    `json:"message,omitempty"`
}

// ProductsListResponse represents paginated API response for product list operations
type ProductsListResponse struct {
	Success    bool           `json:"success"`
	Data       []ProductDTO   `json:"data,omitempty"`
	Error      string         `json:"error,omitempty"`
	Message    string         `json:"message,omitempty"`
	Pagination PaginationInfo `json:"pagination"`
}
