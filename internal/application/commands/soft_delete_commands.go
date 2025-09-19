package commands

import (
	"context"
	"goclean/internal/domain/services"

	"github.com/google/uuid"
)

// SoftDeleteUserCommand represents a command to soft delete a user
type SoftDeleteUserCommand struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	Reason string    `json:"reason,omitempty"` // Optional reason for deletion
}

// RestoreUserCommand represents a command to restore a soft deleted user
type RestoreUserCommand struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	Reason string    `json:"reason,omitempty"` // Optional reason for restoration
}

// SoftDeleteUserCommandHandler handles soft delete user commands
type SoftDeleteUserCommandHandler struct {
	userAggregateService *services.UserAggregateService
}

// NewSoftDeleteUserCommandHandler creates a new soft delete user command handler
func NewSoftDeleteUserCommandHandler(userAggregateService *services.UserAggregateService) *SoftDeleteUserCommandHandler {
	return &SoftDeleteUserCommandHandler{
		userAggregateService: userAggregateService,
	}
}

// Handle handles the soft delete user command
func (h *SoftDeleteUserCommandHandler) Handle(ctx context.Context, cmd SoftDeleteUserCommand) error {
	return h.userAggregateService.SoftDeleteUser(ctx, cmd.UserID)
}

// RestoreUserCommandHandler handles restore user commands
type RestoreUserCommandHandler struct {
	userAggregateService *services.UserAggregateService
}

// NewRestoreUserCommandHandler creates a new restore user command handler
func NewRestoreUserCommandHandler(userAggregateService *services.UserAggregateService) *RestoreUserCommandHandler {
	return &RestoreUserCommandHandler{
		userAggregateService: userAggregateService,
	}
}

// Handle handles the restore user command
func (h *RestoreUserCommandHandler) Handle(ctx context.Context, cmd RestoreUserCommand) error {
	return h.userAggregateService.RestoreUser(ctx, cmd.UserID)
}
