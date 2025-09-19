package queries

import (
	"context"
	"goclean/internal/domain/entities"

	"github.com/google/uuid"
)

// Query represents a query interface
type Query interface{}

// QueryHandler represents a query handler interface
type QueryHandler[T Query, R any] interface {
	Handle(ctx context.Context, query T) (R, error)
}

// GetUserByIDQuery represents a query to get user by ID
type GetUserByIDQuery struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

// GetUserByEmailQuery represents a query to get user by email
type GetUserByEmailQuery struct {
	Email string `json:"email" validate:"required,email"`
}

// GetUserByUsernameQuery represents a query to get user by username
type GetUserByUsernameQuery struct {
	Username string `json:"username" validate:"required"`
}

// ListUsersQuery represents a query to list users
type ListUsersQuery struct {
	Offset int `json:"offset" validate:"min=0"`
	Limit  int `json:"limit" validate:"min=1,max=100"`
}

// GetProductByIDQuery represents a query to get product by ID
type GetProductByIDQuery struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

// GetProductBySKUQuery represents a query to get product by SKU
type GetProductBySKUQuery struct {
	SKU string `json:"sku" validate:"required"`
}

// ListProductsQuery represents a query to list products
type ListProductsQuery struct {
	Offset int `json:"offset" validate:"min=0"`
	Limit  int `json:"limit" validate:"min=1,max=100"`
}

// ListProductsByCategoryQuery represents a query to list products by category
type ListProductsByCategoryQuery struct {
	Category string `json:"category" validate:"required"`
	Offset   int    `json:"offset" validate:"min=0"`
	Limit    int    `json:"limit" validate:"min=1,max=100"`
}

// SearchProductsQuery represents a query to search products
type SearchProductsQuery struct {
	Query  string `json:"query" validate:"required"`
	Offset int    `json:"offset" validate:"min=0"`
	Limit  int    `json:"limit" validate:"min=1,max=100"`
}

// GetOrderByIDQuery represents a query to get order by ID
type GetOrderByIDQuery struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

// GetOrdersByUserIDQuery represents a query to get orders by user ID
type GetOrdersByUserIDQuery struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	Offset int       `json:"offset" validate:"min=0"`
	Limit  int       `json:"limit" validate:"min=1,max=100"`
}

// ListOrdersQuery represents a query to list orders
type ListOrdersQuery struct {
	Offset int `json:"offset" validate:"min=0"`
	Limit  int `json:"limit" validate:"min=1,max=100"`
}

// Query Results

// UserResult represents user query result
type UserResult struct {
	User    *entities.User    `json:"user"`
	Profile *entities.Profile `json:"profile,omitempty"`
}

// UsersResult represents users list query result
type UsersResult struct {
	Users []UserResult `json:"users"`
	Total int          `json:"total"`
}

// ProductResult represents product query result
type ProductResult struct {
	Product *entities.Product `json:"product"`
}

// ProductsResult represents products list query result
type ProductsResult struct {
	Products []*entities.Product `json:"products"`
	Total    int                 `json:"total"`
}

// OrderResult represents order query result
type OrderResult struct {
	Order *entities.Order `json:"order"`
}

// OrdersResult represents orders list query result
type OrdersResult struct {
	Orders []*entities.Order `json:"orders"`
	Total  int               `json:"total"`
}
