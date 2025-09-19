package queries

import (
	"context"
	"goclean/internal/domain/entities"
	"goclean/internal/domain/services"

	"github.com/google/uuid"
)

// GetDeletedUsersQuery represents a query to get deleted users
type GetDeletedUsersQuery struct {
	Offset int `json:"offset" validate:"min=0"`
	Limit  int `json:"limit" validate:"min=1,max=100"`
}

// GetUserWithDeletedQuery represents a query to get a user including deleted ones
type GetUserWithDeletedQuery struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

// DeletedUsersQueryHandler handles queries for deleted users
type DeletedUsersQueryHandler struct {
	userAggregateService *services.UserAggregateService
}

// NewDeletedUsersQueryHandler creates a new deleted users query handler
func NewDeletedUsersQueryHandler(userAggregateService *services.UserAggregateService) *DeletedUsersQueryHandler {
	return &DeletedUsersQueryHandler{
		userAggregateService: userAggregateService,
	}
}

// GetDeletedUsers handles the get deleted users query
func (h *DeletedUsersQueryHandler) GetDeletedUsers(ctx context.Context, query GetDeletedUsersQuery) ([]*entities.User, error) {
	return h.userAggregateService.ListDeletedUsers(ctx, query.Offset, query.Limit)
}

// GetUserWithDeleted handles the get user with deleted query
func (h *DeletedUsersQueryHandler) GetUserWithDeleted(ctx context.Context, query GetUserWithDeletedQuery) (*entities.User, error) {
	return h.userAggregateService.GetUserWithDeleted(ctx, query.UserID)
}
