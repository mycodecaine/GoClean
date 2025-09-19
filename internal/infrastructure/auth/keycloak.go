package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/golang-jwt/jwt/v5"
)

// KeycloakConfig holds Keycloak configuration
type KeycloakConfig struct {
	BaseURL      string
	Realm        string
	ClientID     string
	ClientSecret string
}

// AuthService provides authentication functionality
type AuthService struct {
	client *gocloak.GoCloak
	config KeycloakConfig
}

// NewAuthService creates a new auth service
func NewAuthService(config KeycloakConfig) *AuthService {
	client := gocloak.NewClient(config.BaseURL)

	return &AuthService{
		client: client,
		config: config,
	}
}

// ValidateToken validates JWT token
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*UserClaims, error) {
	// Remove Bearer prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Get public key from Keycloak
	publicKey, err := s.getPublicKey(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}

	// Parse and validate token
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// GetUserInfo retrieves user information from Keycloak
func (s *AuthService) GetUserInfo(ctx context.Context, accessToken string) (*gocloak.UserInfo, error) {
	userInfo, err := s.client.GetUserInfo(ctx, accessToken, s.config.Realm)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	return userInfo, nil
}

// Login authenticates user with Keycloak
func (s *AuthService) Login(ctx context.Context, username, password string) (*gocloak.JWT, error) {
	jwt, err := s.client.Login(ctx, s.config.ClientID, s.config.ClientSecret, s.config.Realm, username, password)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}
	return jwt, nil
}

// RefreshToken refreshes JWT token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error) {
	jwt, err := s.client.RefreshToken(ctx, refreshToken, s.config.ClientID, s.config.ClientSecret, s.config.Realm)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	return jwt, nil
}

// Logout logs out user
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	err := s.client.Logout(ctx, s.config.ClientID, s.config.ClientSecret, s.config.Realm, refreshToken)
	if err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}
	return nil
}

// getPublicKey retrieves public key from Keycloak JWKS endpoint
func (s *AuthService) getPublicKey(ctx context.Context) (interface{}, error) {
	// Get JWKS (JSON Web Key Set) from Keycloak
	jwks, err := s.client.GetCerts(ctx, s.config.Realm)
	if err != nil {
		return nil, fmt.Errorf("failed to get JWKS: %w", err)
	}

	if jwks.Keys == nil || len(*jwks.Keys) == 0 {
		return nil, errors.New("no keys found in JWKS")
	}

	// Use the first key (in production, you might want to find the correct key by kid)
	keys := *jwks.Keys
	key := keys[0]

	// Convert the key to RSA public key
	if key.Kty == nil || *key.Kty != "RSA" {
		return nil, errors.New("unsupported key type, expected RSA")
	}

	// Parse the key - gocloak should provide the key in the proper format
	// For simplicity, we'll use the key directly from the JWKS
	// In a production environment, you might want to implement proper JWK to RSA conversion
	return &key, nil
}

// UserClaims represents JWT claims
type UserClaims struct {
	UserID        string   `json:"sub"`
	Email         string   `json:"email"`
	Username      string   `json:"preferred_username"`
	FirstName     string   `json:"given_name"`
	LastName      string   `json:"family_name"`
	Roles         []string `json:"realm_access.roles"`
	EmailVerified bool     `json:"email_verified"`
	SessionState  string   `json:"session_state"`
	jwt.RegisteredClaims
}

// HasRole checks if user has specific role
func (c *UserClaims) HasRole(role string) bool {
	for _, r := range c.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// IsExpired checks if token is expired
func (c *UserClaims) IsExpired() bool {
	return c.ExpiresAt.Time.Before(time.Now())
}
