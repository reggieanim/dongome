package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dongome/internal/users/app"
	"dongome/internal/users/domain"
	"dongome/internal/users/infra"
	"dongome/pkg/config"
	"dongome/pkg/db"
	"dongome/pkg/events"
	"dongome/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	if err := logger.Initialize(cfg.Server.Mode); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting Dongome API server")

	// Initialize database
	database, err := db.NewDatabase(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", logger.Logger.Error(err))
	}
	defer database.Close()

	// Auto-migrate database schemas
	if err := database.AutoMigrate(
		&domain.User{},
		&domain.SellerProfile{},
	); err != nil {
		logger.Fatal("Failed to migrate database", logger.Logger.Error(err))
	}

	// Initialize NATS event bus
	eventBus, err := events.NewNATSEventBus(cfg.NATS.URL)
	if err != nil {
		logger.Fatal("Failed to connect to NATS", logger.Logger.Error(err))
	}
	defer eventBus.Close()

	// Initialize repositories
	userRepo := infra.NewUserGORMRepository(database.DB)

	// Initialize services
	userService := app.NewUserService(userRepo, eventBus)

	// Initialize handlers
	userHandler := infra.NewUserHandler(userService)

	// Setup Gin router
	if cfg.Server.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
		})
	})

	// API routes
	v1 := router.Group("/api/v1")
	{
		userHandler.RegisterRoutes(v1)
	}

	// Setup server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info(fmt.Sprintf("Server starting on port %s", cfg.Server.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", logger.Logger.Error(err))
		}
	}()

	// Setup event subscriptions
	setupEventSubscriptions(eventBus)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Server shutting down...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", logger.Logger.Error(err))
	}

	logger.Info("Server shutdown complete")
}

// setupEventSubscriptions sets up NATS event subscriptions for cross-bounded context communication
func setupEventSubscriptions(eventBus events.EventBus) {
	// Subscribe to UserRegistered events for notifications
	err := eventBus.Subscribe(domain.UserRegisteredEvent, handleUserRegistered)
	if err != nil {
		logger.Error("Failed to subscribe to UserRegistered events", logger.Logger.Error(err))
	}

	// Subscribe to UserEmailVerified events
	err = eventBus.Subscribe(domain.UserEmailVerifiedEvent, handleUserEmailVerified)
	if err != nil {
		logger.Error("Failed to subscribe to UserEmailVerified events", logger.Logger.Error(err))
	}

	logger.Info("Event subscriptions setup complete")
}

// Event handlers for demonstration of cross-bounded context communication
func handleUserRegistered(ctx context.Context, event *events.Event) error {
	logger.Info("Handling UserRegistered event",
		logger.Logger.String("event_id", event.ID),
		logger.Logger.String("user_id", event.AggregateID))

	var userData domain.UserRegistered
	if err := events.ParseEventData(event, &userData); err != nil {
		return err
	}

	// In a real application, this would:
	// 1. Send welcome email to the user
	// 2. Create user profile in other services
	// 3. Send verification email
	// 4. Add to analytics/metrics

	logger.Info("UserRegistered event processed successfully",
		logger.Logger.String("user_email", userData.Email),
		logger.Logger.String("user_name", userData.FirstName+" "+userData.LastName))

	// Simulate sending welcome email
	// emailService.SendWelcomeEmail(userData.Email, userData.VerificationToken)

	return nil
}

func handleUserEmailVerified(ctx context.Context, event *events.Event) error {
	logger.Info("Handling UserEmailVerified event",
		logger.Logger.String("event_id", event.ID),
		logger.Logger.String("user_id", event.AggregateID))

	var userData domain.UserEmailVerified
	if err := events.ParseEventData(event, &userData); err != nil {
		return err
	}

	// In a real application, this would:
	// 1. Update user status in other services
	// 2. Send confirmation email
	// 3. Enable additional features
	// 4. Update analytics

	logger.Info("UserEmailVerified event processed successfully",
		logger.Logger.String("user_email", userData.Email))

	return nil
}
