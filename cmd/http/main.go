package main

import (
	"context"
	"goclean/internal/application/commands"
	"goclean/internal/application/queries"
	"goclean/internal/domain/services"
	"goclean/internal/infrastructure/auth"
	"goclean/internal/infrastructure/cache"
	"goclean/internal/infrastructure/persistence"
	gormPersistence "goclean/internal/infrastructure/persistence/gorm"
	httpServer "goclean/internal/interfaces/http"
	"goclean/internal/interfaces/http/handlers"
	"goclean/pkg/config"
	"goclean/pkg/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title GoClean API
// @version 1.0
// @description A clean architecture Go API with DDD, CQRS, REST, and gRPC
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	loggerConfig := logger.Config{
		Level:  logger.LogLevel(cfg.App.LogLevel),
		Format: "text",
		Output: "stdout",
	}
	if cfg.IsProduction() {
		loggerConfig.Format = "json"
	}

	appLogger, err := logger.New(loggerConfig)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	appLogger.Info("Starting GoClean API",
		"version", cfg.App.Version,
		"environment", cfg.App.Environment,
	)

	// Initialize database
	db, err := persistence.NewDatabase(persistence.DatabaseConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	})
	if err != nil {
		appLogger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	// Run database migrations
	if err := persistence.MigrateDatabase(db); err != nil {
		appLogger.Error("Failed to migrate database", "error", err)
		os.Exit(1)
	}
	appLogger.Info("Database migrations completed successfully")

	// Initialize cache service
	cacheService := cache.NewCacheService(cache.RedisConfig{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test cache connection
	if err := cacheService.Ping(context.Background()); err != nil {
		appLogger.Error("Failed to connect to Redis", "error", err)
		os.Exit(1)
	}
	appLogger.Info("Connected to Redis successfully")

	// Initialize auth service
	authService := auth.NewAuthService(auth.KeycloakConfig{
		BaseURL:      cfg.Keycloak.BaseURL,
		Realm:        cfg.Keycloak.Realm,
		ClientID:     cfg.Keycloak.ClientID,
		ClientSecret: cfg.Keycloak.ClientSecret,
	})

	// Initialize repositories
	userRepo := gormPersistence.NewUserRepository(db)
	profileRepo := persistence.NewProfileGormRepository(db)
	productRepo := persistence.NewProductGormRepository(db)
	orderRepo := persistence.NewOrderGormRepository(db)

	// Initialize domain services
	userDomainService := services.NewUserDomainService(userRepo, profileRepo)
	productDomainService := services.NewProductDomainService(productRepo)
	orderDomainService := services.NewOrderDomainService(orderRepo, productRepo)

	// Initialize command handlers
	userCommandHandler := commands.NewUserCommandHandler(userDomainService)
	productCommandHandler := commands.NewProductCommandHandler(productDomainService)
	orderCommandHandler := commands.NewOrderCommandHandler(orderDomainService)

	// Initialize query handlers
	userQueryHandler := queries.NewUserQueryHandler(userRepo, profileRepo)
	productQueryHandler := queries.NewProductQueryHandler(productRepo)
	orderQueryHandler := queries.NewOrderQueryHandler(orderRepo)

	// Initialize HTTP handlers
	userHandler := handlers.NewUserHandler(userCommandHandler, userQueryHandler)
	productHandler := handlers.NewProductHandler(productCommandHandler, productQueryHandler)

	// Note: Order handlers are available but not yet integrated into the server
	_ = orderCommandHandler // TODO: Integrate order handlers
	_ = orderQueryHandler   // TODO: Integrate order handlers

	// Initialize HTTP server
	server := httpServer.NewServer(
		appLogger,
		authService,
		userHandler,
		productHandler,
	)

	// Start HTTP server in a goroutine
	go func() {
		if err := server.Start(cfg.GetServerAddress()); err != nil && err != http.ErrServerClosed {
			appLogger.Error("Failed to start HTTP server", "error", err)
			os.Exit(1)
		}
	}()

	appLogger.Info("HTTP server started", "address", cfg.GetServerAddress())

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	// Create a deadline to wait for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server gracefully
	if err := server.Echo().Shutdown(ctx); err != nil {
		appLogger.Error("Failed to gracefully shutdown HTTP server", "error", err)
	}

	// Close cache connection
	if err := cacheService.Close(); err != nil {
		appLogger.Error("Failed to close cache connection", "error", err)
	}

	// Close database connection
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}

	appLogger.Info("Server exited")
}
