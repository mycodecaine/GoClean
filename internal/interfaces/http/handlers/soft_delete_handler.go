package handlers

import (
	"net/http"
	"strconv"

	"goclean/internal/application/commands"
	"goclean/internal/application/queries"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// SoftDeleteHandler handles soft delete operations for various entities
type SoftDeleteHandler struct {
	softDeleteUserHandler *commands.SoftDeleteUserCommandHandler
	restoreUserHandler    *commands.RestoreUserCommandHandler
	deletedUsersHandler   *queries.DeletedUsersQueryHandler
}

// NewSoftDeleteHandler creates a new soft delete handler
func NewSoftDeleteHandler(
	softDeleteUserHandler *commands.SoftDeleteUserCommandHandler,
	restoreUserHandler *commands.RestoreUserCommandHandler,
	deletedUsersHandler *queries.DeletedUsersQueryHandler,
) *SoftDeleteHandler {
	return &SoftDeleteHandler{
		softDeleteUserHandler: softDeleteUserHandler,
		restoreUserHandler:    restoreUserHandler,
		deletedUsersHandler:   deletedUsersHandler,
	}
}

// SoftDeleteUser godoc
// @Summary Soft delete a user
// @Description Marks a user as deleted without permanently removing from database
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param body body commands.SoftDeleteUserCommand true "Soft delete command"
// @Success 204 "User soft deleted successfully"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/soft-delete [delete]
func (h *SoftDeleteHandler) SoftDeleteUser(c echo.Context) error {
	idStr := c.Param("id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	var req commands.SoftDeleteUserCommand
	req.UserID = userID

	if err := c.Bind(&req); err == nil {
		// If body is provided, use it (for reason, etc.)
	}

	if err := h.softDeleteUserHandler.Handle(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// RestoreUser godoc
// @Summary Restore a soft deleted user
// @Description Restores a previously soft deleted user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param body body commands.RestoreUserCommand true "Restore command"
// @Success 204 "User restored successfully"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/restore [post]
func (h *SoftDeleteHandler) RestoreUser(c echo.Context) error {
	idStr := c.Param("id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	var req commands.RestoreUserCommand
	req.UserID = userID

	if err := c.Bind(&req); err == nil {
		// If body is provided, use it (for reason, etc.)
	}

	if err := h.restoreUserHandler.Handle(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetDeletedUsers godoc
// @Summary Get deleted users
// @Description Retrieves a list of soft deleted users
// @Tags users
// @Accept json
// @Produce json
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Success 200 {array} entities.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/deleted [get]
func (h *SoftDeleteHandler) GetDeletedUsers(c echo.Context) error {
	offsetStr := c.QueryParam("offset")
	limitStr := c.QueryParam("limit")

	offset := 0
	limit := 10

	if offsetStr != "" {
		if parsed, err := strconv.Atoi(offsetStr); err == nil {
			offset = parsed
		}
	}

	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil {
			limit = parsed
		}
	}

	query := queries.GetDeletedUsersQuery{
		Offset: offset,
		Limit:  limit,
	}

	users, err := h.deletedUsersHandler.GetDeletedUsers(c.Request().Context(), query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, users)
}

// GetUserWithDeleted godoc
// @Summary Get user including deleted
// @Description Retrieves a user by ID including soft deleted users
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} entities.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id}/with-deleted [get]
func (h *SoftDeleteHandler) GetUserWithDeleted(c echo.Context) error {
	idStr := c.Param("id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	query := queries.GetUserWithDeletedQuery{
		UserID: userID,
	}

	user, err := h.deletedUsersHandler.GetUserWithDeleted(c.Request().Context(), query)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}
