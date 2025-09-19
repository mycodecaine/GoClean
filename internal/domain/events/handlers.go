package events

import (
	"context"
	"goclean/internal/domain/entities"
	"goclean/pkg/logger"
)

// UserCreatedEventHandler handles user created events
type UserCreatedEventHandler struct {
	logger *logger.Logger
}

// NewUserCreatedEventHandler creates a new user created event handler
func NewUserCreatedEventHandler(logger *logger.Logger) *UserCreatedEventHandler {
	return &UserCreatedEventHandler{
		logger: logger,
	}
}

// Handle handles the user created event
func (h *UserCreatedEventHandler) Handle(ctx context.Context, event entities.DomainEvent) error {
	userCreatedEvent, ok := event.(entities.UserCreatedEvent)
	if !ok {
		return nil
	}

	h.logger.Info("User created event handled",
		"user_id", userCreatedEvent.UserID,
		"email", userCreatedEvent.Email,
		"username", userCreatedEvent.Username,
	)

	// Here you can add additional logic like:
	// - Send welcome email
	// - Create default user settings
	// - Update analytics
	// - Publish to message queue

	return nil
}

// CanHandle checks if this handler can handle the event
func (h *UserCreatedEventHandler) CanHandle(event entities.DomainEvent) bool {
	_, ok := event.(entities.UserCreatedEvent)
	return ok
}

// UserDeletedEventHandler handles user deleted events
type UserDeletedEventHandler struct {
	logger *logger.Logger
}

// NewUserDeletedEventHandler creates a new user deleted event handler
func NewUserDeletedEventHandler(logger *logger.Logger) *UserDeletedEventHandler {
	return &UserDeletedEventHandler{
		logger: logger,
	}
}

// Handle handles the user deleted event
func (h *UserDeletedEventHandler) Handle(ctx context.Context, event entities.DomainEvent) error {
	userDeletedEvent, ok := event.(entities.UserDeletedEvent)
	if !ok {
		return nil
	}

	h.logger.Info("User deleted event handled",
		"user_id", userDeletedEvent.UserID,
	)

	// Here you can add additional logic like:
	// - Archive user data
	// - Cancel active subscriptions
	// - Clean up external references
	// - Send deletion confirmation

	return nil
}

// CanHandle checks if this handler can handle the event
func (h *UserDeletedEventHandler) CanHandle(event entities.DomainEvent) bool {
	_, ok := event.(entities.UserDeletedEvent)
	return ok
}

// ProductCreatedEventHandler handles product created events
type ProductCreatedEventHandler struct {
	logger *logger.Logger
}

// NewProductCreatedEventHandler creates a new product created event handler
func NewProductCreatedEventHandler(logger *logger.Logger) *ProductCreatedEventHandler {
	return &ProductCreatedEventHandler{
		logger: logger,
	}
}

// Handle handles the product created event
func (h *ProductCreatedEventHandler) Handle(ctx context.Context, event entities.DomainEvent) error {
	productCreatedEvent, ok := event.(entities.ProductCreatedEvent)
	if !ok {
		return nil
	}

	h.logger.Info("Product created event handled",
		"product_id", productCreatedEvent.ProductID,
		"name", productCreatedEvent.Name,
		"sku", productCreatedEvent.SKU,
	)

	// Here you can add additional logic like:
	// - Update search index
	// - Generate product recommendations
	// - Update inventory system
	// - Send notifications to subscribers

	return nil
}

// CanHandle checks if this handler can handle the event
func (h *ProductCreatedEventHandler) CanHandle(event entities.DomainEvent) bool {
	_, ok := event.(entities.ProductCreatedEvent)
	return ok
}
