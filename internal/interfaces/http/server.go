package http

import (
	"goclean/internal/infrastructure/auth"
	"goclean/internal/interfaces/http/handlers"
	"goclean/internal/interfaces/http/middleware"
	"goclean/pkg/logger"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// Server represents the HTTP server
type Server struct {
	echo        *echo.Echo
	logger      *logger.Logger
	authService *auth.AuthService
}

// NewServer creates a new HTTP server
func NewServer(
	logger *logger.Logger,
	authService *auth.AuthService,
	userHandler *handlers.UserHandler,
	productHandler *handlers.ProductHandler,
) *Server {
	e := echo.New()

	// Configure Echo
	e.HideBanner = true
	e.HidePort = true

	server := &Server{
		echo:        e,
		logger:      logger,
		authService: authService,
	}

	// Setup middleware
	server.setupMiddleware()

	// Setup routes
	server.setupRoutes(userHandler, productHandler)

	return server
}

// setupMiddleware configures middleware
func (s *Server) setupMiddleware() {
	// Basic middleware
	s.echo.Use(echoMiddleware.Logger())
	s.echo.Use(echoMiddleware.Recover())
	s.echo.Use(echoMiddleware.CORS())
	s.echo.Use(echoMiddleware.Secure())
	s.echo.Use(echoMiddleware.RequestID())

	// Custom request logging
	s.echo.Use(echoMiddleware.RequestLoggerWithConfig(echoMiddleware.RequestLoggerConfig{
		LogURI:     true,
		LogStatus:  true,
		LogMethod:  true,
		LogLatency: true,
		LogValuesFunc: func(c echo.Context, v echoMiddleware.RequestLoggerValues) error {
			s.logger.Info("HTTP request",
				"method", v.Method,
				"uri", v.URI,
				"status", v.Status,
				"latency", v.Latency,
			)
			return nil
		},
	}))
}

// setupRoutes configures API routes
func (s *Server) setupRoutes(
	userHandler *handlers.UserHandler,
	productHandler *handlers.ProductHandler,
) {
	// Health check
	s.echo.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status":  "healthy",
			"service": "goclean-api",
		})
	})

	// Swagger documentation
	s.echo.GET("/swagger/*", echoSwagger.WrapHandler)

	// Auth middleware
	authMiddleware := middleware.NewAuthMiddleware(s.authService, []string{
		"/health",
		"/swagger",
		"/api/v1/auth",
	})

	// API v1 routes
	api := s.echo.Group("/api/v1")

	// Public routes (no auth required)
	public := api.Group("")

	// Protected routes (auth required)
	protected := api.Group("")
	protected.Use(authMiddleware.Authenticate)

	// User routes
	public.POST("/users", userHandler.CreateUser)          // Public registration
	protected.GET("/users", userHandler.ListUsers)         // Admin only
	protected.GET("/users/:id", userHandler.GetUser)       // Auth required
	protected.GET("/users/me", userHandler.GetCurrentUser) // Auth required

	// Product routes
	public.GET("/products", productHandler.ListProducts)      // Public
	public.GET("/products/:id", productHandler.GetProduct)    // Public
	protected.POST("/products", productHandler.CreateProduct) // Auth required

	// Admin routes (require admin role)
	admin := protected.Group("/admin")
	admin.Use(authMiddleware.RequireRole("admin"))

	// Additional admin routes can be added here
}

// Start starts the HTTP server
func (s *Server) Start(address string) error {
	s.logger.Info("Starting HTTP server", "address", address)
	return s.echo.Start(address)
}

// Stop stops the HTTP server gracefully
func (s *Server) Stop() error {
	s.logger.Info("Stopping HTTP server")
	return s.echo.Close()
}

// Echo returns the underlying Echo instance
func (s *Server) Echo() *echo.Echo {
	return s.echo
}
