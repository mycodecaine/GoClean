package queries

import (
	"context"
	"goclean/internal/domain/repositories"
)

// UserQueryHandler handles user-related queries
type UserQueryHandler struct {
	userRepo    repositories.UserRepository
	profileRepo repositories.ProfileRepository
}

// NewUserQueryHandler creates a new user query handler
func NewUserQueryHandler(userRepo repositories.UserRepository, profileRepo repositories.ProfileRepository) *UserQueryHandler {
	return &UserQueryHandler{
		userRepo:    userRepo,
		profileRepo: profileRepo,
	}
}

// Handle handles GetUserByIDQuery
func (h *UserQueryHandler) Handle(ctx context.Context, query GetUserByIDQuery) (*UserResult, error) {
	user, err := h.userRepo.GetByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}

	profile, _ := h.profileRepo.GetByUserID(ctx, user.ID)

	return &UserResult{
		User:    user,
		Profile: profile,
	}, nil
}

// HandleByEmail handles GetUserByEmailQuery
func (h *UserQueryHandler) HandleByEmail(ctx context.Context, query GetUserByEmailQuery) (*UserResult, error) {
	user, err := h.userRepo.GetByEmail(ctx, query.Email)
	if err != nil {
		return nil, err
	}

	profile, _ := h.profileRepo.GetByUserID(ctx, user.ID)

	return &UserResult{
		User:    user,
		Profile: profile,
	}, nil
}

// HandleByUsername handles GetUserByUsernameQuery
func (h *UserQueryHandler) HandleByUsername(ctx context.Context, query GetUserByUsernameQuery) (*UserResult, error) {
	user, err := h.userRepo.GetByUsername(ctx, query.Username)
	if err != nil {
		return nil, err
	}

	profile, _ := h.profileRepo.GetByUserID(ctx, user.ID)

	return &UserResult{
		User:    user,
		Profile: profile,
	}, nil
}

// HandleList handles ListUsersQuery
func (h *UserQueryHandler) HandleList(ctx context.Context, query ListUsersQuery) (*UsersResult, error) {
	users, err := h.userRepo.List(ctx, query.Offset, query.Limit)
	if err != nil {
		return nil, err
	}

	results := make([]UserResult, len(users))
	for i, user := range users {
		profile, _ := h.profileRepo.GetByUserID(ctx, user.ID)
		results[i] = UserResult{
			User:    user,
			Profile: profile,
		}
	}

	return &UsersResult{
		Users: results,
		Total: len(results), // In a real implementation, you'd get total count separately
	}, nil
}

// ProductQueryHandler handles product-related queries
type ProductQueryHandler struct {
	productRepo repositories.ProductRepository
}

// NewProductQueryHandler creates a new product query handler
func NewProductQueryHandler(productRepo repositories.ProductRepository) *ProductQueryHandler {
	return &ProductQueryHandler{
		productRepo: productRepo,
	}
}

// Handle handles GetProductByIDQuery
func (h *ProductQueryHandler) Handle(ctx context.Context, query GetProductByIDQuery) (*ProductResult, error) {
	product, err := h.productRepo.GetByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}

	return &ProductResult{Product: product}, nil
}

// HandleBySKU handles GetProductBySKUQuery
func (h *ProductQueryHandler) HandleBySKU(ctx context.Context, query GetProductBySKUQuery) (*ProductResult, error) {
	product, err := h.productRepo.GetBySKU(ctx, query.SKU)
	if err != nil {
		return nil, err
	}

	return &ProductResult{Product: product}, nil
}

// HandleList handles ListProductsQuery
func (h *ProductQueryHandler) HandleList(ctx context.Context, query ListProductsQuery) (*ProductsResult, error) {
	products, err := h.productRepo.List(ctx, query.Offset, query.Limit)
	if err != nil {
		return nil, err
	}

	return &ProductsResult{
		Products: products,
		Total:    len(products), // In a real implementation, you'd get total count separately
	}, nil
}

// HandleByCategory handles ListProductsByCategoryQuery
func (h *ProductQueryHandler) HandleByCategory(ctx context.Context, query ListProductsByCategoryQuery) (*ProductsResult, error) {
	products, err := h.productRepo.ListByCategory(ctx, query.Category, query.Offset, query.Limit)
	if err != nil {
		return nil, err
	}

	return &ProductsResult{
		Products: products,
		Total:    len(products), // In a real implementation, you'd get total count separately
	}, nil
}

// HandleSearch handles SearchProductsQuery
func (h *ProductQueryHandler) HandleSearch(ctx context.Context, query SearchProductsQuery) (*ProductsResult, error) {
	products, err := h.productRepo.Search(ctx, query.Query, query.Offset, query.Limit)
	if err != nil {
		return nil, err
	}

	return &ProductsResult{
		Products: products,
		Total:    len(products), // In a real implementation, you'd get total count separately
	}, nil
}

// OrderQueryHandler handles order-related queries
type OrderQueryHandler struct {
	orderRepo repositories.OrderRepository
}

// NewOrderQueryHandler creates a new order query handler
func NewOrderQueryHandler(orderRepo repositories.OrderRepository) *OrderQueryHandler {
	return &OrderQueryHandler{
		orderRepo: orderRepo,
	}
}

// Handle handles GetOrderByIDQuery
func (h *OrderQueryHandler) Handle(ctx context.Context, query GetOrderByIDQuery) (*OrderResult, error) {
	order, err := h.orderRepo.GetByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}

	return &OrderResult{Order: order}, nil
}

// HandleByUserID handles GetOrdersByUserIDQuery
func (h *OrderQueryHandler) HandleByUserID(ctx context.Context, query GetOrdersByUserIDQuery) (*OrdersResult, error) {
	orders, err := h.orderRepo.GetByUserID(ctx, query.UserID, query.Offset, query.Limit)
	if err != nil {
		return nil, err
	}

	return &OrdersResult{
		Orders: orders,
		Total:  len(orders), // In a real implementation, you'd get total count separately
	}, nil
}
