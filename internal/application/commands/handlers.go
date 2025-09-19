package commands

import (
	"context"
	"goclean/internal/domain/entities"
	"goclean/internal/domain/services"
	"time"

	"github.com/google/uuid"
)

// UserCommandHandler handles user-related commands
type UserCommandHandler struct {
	userService *services.UserDomainService
}

// NewUserCommandHandler creates a new user command handler
func NewUserCommandHandler(userService *services.UserDomainService) *UserCommandHandler {
	return &UserCommandHandler{
		userService: userService,
	}
}

// Handle handles CreateUserCommand
func (h *UserCommandHandler) Handle(ctx context.Context, cmd CreateUserCommand) error {
	user := entities.NewUser(cmd.Email, cmd.Username, cmd.FirstName, cmd.LastName)

	var profile *entities.Profile
	if cmd.Profile != nil {
		profile = entities.NewProfile(user.ID, cmd.Profile.Bio, cmd.Profile.Avatar, nil)

		if cmd.Profile.DateOfBirth != "" {
			if dateOfBirth, err := time.Parse("2006-01-02", cmd.Profile.DateOfBirth); err == nil {
				profile.DateOfBirth = &dateOfBirth
			}
		}
	}

	return h.userService.CreateUserWithProfile(ctx, user, profile)
}

// ProductCommandHandler handles product-related commands
type ProductCommandHandler struct {
	productService *services.ProductDomainService
}

// NewProductCommandHandler creates a new product command handler
func NewProductCommandHandler(productService *services.ProductDomainService) *ProductCommandHandler {
	return &ProductCommandHandler{
		productService: productService,
	}
}

// Handle handles CreateProductCommand
func (h *ProductCommandHandler) Handle(ctx context.Context, cmd CreateProductCommand) error {
	product := entities.NewProduct(cmd.Name, cmd.Description, cmd.SKU, cmd.Category, cmd.Price, cmd.CreatedBy)

	if err := h.productService.ValidateProduct(product); err != nil {
		return err
	}

	return h.productService.CreateProduct(ctx, product)
}

// OrderCommandHandler handles order-related commands
type OrderCommandHandler struct {
	orderService *services.OrderDomainService
}

// NewOrderCommandHandler creates a new order command handler
func NewOrderCommandHandler(orderService *services.OrderDomainService) *OrderCommandHandler {
	return &OrderCommandHandler{
		orderService: orderService,
	}
}

// Handle handles CreateOrderCommand
func (h *OrderCommandHandler) Handle(ctx context.Context, cmd CreateOrderCommand) error {
	items := make([]entities.OrderItem, len(cmd.Items))
	for i, item := range cmd.Items {
		// Note: Price would typically be fetched from ProductRepository
		items[i] = *entities.NewOrderItem(uuid.Nil, item.ProductID, item.Quantity, 0.0) // OrderID set later
	}

	order := entities.NewOrder(cmd.UserID, items)
	return h.orderService.CreateOrder(ctx, order)
}

// HandleUpdateOrderStatus handles UpdateOrderStatusCommand
func (h *OrderCommandHandler) HandleUpdateOrderStatus(ctx context.Context, cmd UpdateOrderStatusCommand) error {
	return h.orderService.UpdateOrderStatus(ctx, cmd.ID, cmd.Status)
}
