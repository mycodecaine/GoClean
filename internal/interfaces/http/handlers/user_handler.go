package handlers

import (
	"context"
	"goclean/internal/application/commands"
	"goclean/internal/application/dto"
	"goclean/internal/application/queries"
	"goclean/internal/infrastructure/auth"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userCommandHandler *commands.UserCommandHandler
	userQueryHandler   *queries.UserQueryHandler
}

// NewUserHandler creates a new user handler
func NewUserHandler(
	userCommandHandler *commands.UserCommandHandler,
	userQueryHandler *queries.UserQueryHandler,
) *UserHandler {
	return &UserHandler{
		userCommandHandler: userCommandHandler,
		userQueryHandler:   userQueryHandler,
	}
}

// CreateUser creates a new user
// @Summary Create a new user
// @Description Create a new user with profile information
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User data"
// @Success 201 {object} dto.UserAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 500 {object} dto.ErrorAPIResponse
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c echo.Context) error {
	var req dto.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse[interface{}]{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Convert to command
	cmd := commands.CreateUserCommand{
		Email:     req.Email,
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if req.Profile != nil {
		cmd.Profile = &commands.CreateProfileData{
			Bio:         req.Profile.Bio,
			Avatar:      req.Profile.Avatar,
			DateOfBirth: req.Profile.DateOfBirth,
		}
	}

	// Execute command
	if err := h.userCommandHandler.Handle(context.Background(), cmd); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.APIResponse[interface{}]{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.APIResponse[interface{}]{
		Success: true,
		Message: "User created successfully",
	})
}

// GetUser retrieves a user by ID
// @Summary Get user by ID
// @Description Get user information by user ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.UserAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 404 {object} dto.ErrorAPIResponse
// @Router /api/v1/users/{id} [get]
// @Security BearerAuth
func (h *UserHandler) GetUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse[interface{}]{
			Success: false,
			Error:   "Invalid user ID",
		})
	}

	query := queries.GetUserByIDQuery{ID: id}
	result, err := h.userQueryHandler.Handle(context.Background(), query)
	if err != nil {
		return c.JSON(http.StatusNotFound, dto.APIResponse[interface{}]{
			Success: false,
			Error:   "User not found",
		})
	}

	// Convert to DTO
	userDTO := &dto.UserDTO{
		ID:        result.User.ID,
		Email:     result.User.Email,
		Username:  result.User.Username,
		FirstName: result.User.FirstName,
		LastName:  result.User.LastName,
		IsActive:  result.User.IsActive,
		CreatedAt: result.User.CreatedAt,
		UpdatedAt: result.User.UpdatedAt,
	}

	if result.Profile != nil {
		userDTO.Profile = &dto.ProfileDTO{
			ID:          result.Profile.ID,
			UserID:      result.Profile.UserID,
			Bio:         result.Profile.Bio,
			Avatar:      result.Profile.Avatar,
			DateOfBirth: result.Profile.DateOfBirth,
			CreatedAt:   result.Profile.CreatedAt,
			UpdatedAt:   result.Profile.UpdatedAt,
		}
	}

	return c.JSON(http.StatusOK, dto.APIResponse[*dto.UserDTO]{
		Success: true,
		Data:    userDTO,
	})
}

// ListUsers retrieves users with pagination
// @Summary List users
// @Description Get list of users with pagination
// @Tags users
// @Produce json
// @Param offset query int false "Offset" default(0)
// @Param limit query int false "Limit" default(10)
// @Success 200 {object} dto.UsersListResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Router /api/v1/users [get]
// @Security BearerAuth
func (h *UserHandler) ListUsers(c echo.Context) error {
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 10
	}

	query := queries.ListUsersQuery{
		Offset: offset,
		Limit:  limit,
	}

	result, err := h.userQueryHandler.HandleList(context.Background(), query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.APIResponse[interface{}]{
			Success: false,
			Error:   err.Error(),
		})
	}

	// Convert to DTOs
	userDTOs := make([]dto.UserDTO, len(result.Users))
	for i, userResult := range result.Users {
		userDTOs[i] = dto.UserDTO{
			ID:        userResult.User.ID,
			Email:     userResult.User.Email,
			Username:  userResult.User.Username,
			FirstName: userResult.User.FirstName,
			LastName:  userResult.User.LastName,
			IsActive:  userResult.User.IsActive,
			CreatedAt: userResult.User.CreatedAt,
			UpdatedAt: userResult.User.UpdatedAt,
		}

		if userResult.Profile != nil {
			userDTOs[i].Profile = &dto.ProfileDTO{
				ID:          userResult.Profile.ID,
				UserID:      userResult.Profile.UserID,
				Bio:         userResult.Profile.Bio,
				Avatar:      userResult.Profile.Avatar,
				DateOfBirth: userResult.Profile.DateOfBirth,
				CreatedAt:   userResult.Profile.CreatedAt,
				UpdatedAt:   userResult.Profile.UpdatedAt,
			}
		}
	}

	return c.JSON(http.StatusOK, dto.PaginatedResponse[[]dto.UserDTO]{
		APIResponse: dto.APIResponse[[]dto.UserDTO]{
			Success: true,
			Data:    userDTOs,
		},
		Pagination: dto.PaginationInfo{
			Offset: offset,
			Limit:  limit,
			Total:  result.Total,
		},
	})
}

// GetCurrentUser retrieves current authenticated user
// @Summary Get current user
// @Description Get current authenticated user information
// @Tags users
// @Produce json
// @Success 200 {object} dto.UserAPIResponse
// @Failure 401 {object} dto.ErrorAPIResponse
// @Failure 404 {object} dto.ErrorAPIResponse
// @Router /api/v1/users/me [get]
// @Security BearerAuth
func (h *UserHandler) GetCurrentUser(c echo.Context) error {
	claims, ok := c.Get("user_claims").(*auth.UserClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.APIResponse[interface{}]{
			Success: false,
			Error:   "User not authenticated",
		})
	}

	id, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse[interface{}]{
			Success: false,
			Error:   "Invalid user ID",
		})
	}

	query := queries.GetUserByIDQuery{ID: id}
	result, err := h.userQueryHandler.Handle(context.Background(), query)
	if err != nil {
		return c.JSON(http.StatusNotFound, dto.APIResponse[interface{}]{
			Success: false,
			Error:   "User not found",
		})
	}

	// Convert to DTO
	userDTO := &dto.UserDTO{
		ID:        result.User.ID,
		Email:     result.User.Email,
		Username:  result.User.Username,
		FirstName: result.User.FirstName,
		LastName:  result.User.LastName,
		IsActive:  result.User.IsActive,
		CreatedAt: result.User.CreatedAt,
		UpdatedAt: result.User.UpdatedAt,
	}

	if result.Profile != nil {
		userDTO.Profile = &dto.ProfileDTO{
			ID:          result.Profile.ID,
			UserID:      result.Profile.UserID,
			Bio:         result.Profile.Bio,
			Avatar:      result.Profile.Avatar,
			DateOfBirth: result.Profile.DateOfBirth,
			CreatedAt:   result.Profile.CreatedAt,
			UpdatedAt:   result.Profile.UpdatedAt,
		}
	}

	return c.JSON(http.StatusOK, dto.APIResponse[*dto.UserDTO]{
		Success: true,
		Data:    userDTO,
	})
}
