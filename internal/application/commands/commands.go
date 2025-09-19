package commands

import (
	"context"
	"goclean/internal/domain/entities"

	"github.com/google/uuid"
)

// Command represents a command interface
type Command interface{}

// CommandHandler represents a command handler interface
type CommandHandler[T Command] interface {
	Handle(ctx context.Context, cmd T) error
}

// CreateUserCommand represents a command to create a user
type CreateUserCommand struct {
	Email     string             `json:"email" validate:"required,email"`
	Username  string             `json:"username" validate:"required,min=3,max=50"`
	FirstName string             `json:"first_name" validate:"required,min=1,max=100"`
	LastName  string             `json:"last_name" validate:"required,min=1,max=100"`
	Profile   *CreateProfileData `json:"profile,omitempty"`
}

type CreateProfileData struct {
	Bio         string `json:"bio"`
	Avatar      string `json:"avatar"`
	DateOfBirth string `json:"date_of_birth"`
}

// UpdateUserCommand represents a command to update a user
type UpdateUserCommand struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	Email     string    `json:"email" validate:"email"`
	Username  string    `json:"username" validate:"min=3,max=50"`
	FirstName string    `json:"first_name" validate:"min=1,max=100"`
	LastName  string    `json:"last_name" validate:"min=1,max=100"`
}

// DeleteUserCommand represents a command to delete a user
type DeleteUserCommand struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

// CreateProductCommand represents a command to create a product
type CreateProductCommand struct {
	Name        string    `json:"name" validate:"required,min=1,max=255"`
	Description string    `json:"description"`
	Price       float64   `json:"price" validate:"required,gt=0"`
	SKU         string    `json:"sku" validate:"required,min=1,max=100"`
	Category    string    `json:"category" validate:"required"`
	CreatedBy   uuid.UUID `json:"created_by" validate:"required"`
}

// UpdateProductCommand represents a command to update a product
type UpdateProductCommand struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Name        string    `json:"name" validate:"min=1,max=255"`
	Description string    `json:"description"`
	Price       float64   `json:"price" validate:"gt=0"`
	Category    string    `json:"category"`
	IsActive    *bool     `json:"is_active"`
}

// DeleteProductCommand represents a command to delete a product
type DeleteProductCommand struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

// CreateOrderCommand represents a command to create an order
type CreateOrderCommand struct {
	UserID uuid.UUID             `json:"user_id" validate:"required"`
	Items  []CreateOrderItemData `json:"items" validate:"required,min=1"`
}

type CreateOrderItemData struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,gt=0"`
}

// UpdateOrderStatusCommand represents a command to update order status
type UpdateOrderStatusCommand struct {
	ID     uuid.UUID            `json:"id" validate:"required"`
	Status entities.OrderStatus `json:"status" validate:"required"`
}

// CancelOrderCommand represents a command to cancel an order
type CancelOrderCommand struct {
	ID uuid.UUID `json:"id" validate:"required"`
}
