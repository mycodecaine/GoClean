package middleware

import (
	"context"
	"goclean/internal/infrastructure/auth"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// AuthMiddleware provides JWT authentication middleware
type AuthMiddleware struct {
	authService *auth.AuthService
	skipPaths   []string
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(authService *auth.AuthService, skipPaths []string) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		skipPaths:   skipPaths,
	}
}

// Authenticate validates JWT token
func (m *AuthMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Check if path should be skipped
		path := c.Request().URL.Path
		for _, skipPath := range m.skipPaths {
			if strings.HasPrefix(path, skipPath) {
				return next(c)
			}
		}

		// Get Authorization header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing authorization header")
		}

		// Validate token
		claims, err := m.authService.ValidateToken(context.Background(), authHeader)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token: "+err.Error())
		}

		// Add user info to context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_username", claims.Username)
		c.Set("user_roles", claims.Roles)
		c.Set("user_claims", claims)

		return next(c)
	}
}

// RequireRole checks if user has required role
func (m *AuthMiddleware) RequireRole(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get("user_claims").(*auth.UserClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
			}

			if !claims.HasRole(role) {
				return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions")
			}

			return next(c)
		}
	}
}

// RequireAnyRole checks if user has any of the required roles
func (m *AuthMiddleware) RequireAnyRole(roles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get("user_claims").(*auth.UserClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
			}

			hasRole := false
			for _, role := range roles {
				if claims.HasRole(role) {
					hasRole = true
					break
				}
			}

			if !hasRole {
				return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions")
			}

			return next(c)
		}
	}
}
