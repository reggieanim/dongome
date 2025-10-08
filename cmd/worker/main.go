package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"dongome/internal/users/domain"
	"dongome/pkg/config"
	"dongome/pkg/events"
	"dongome/pkg/logger"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	if err := logger.Initialize(cfg.Server.Mode); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting Dongome Worker")

	// Initialize NATS event bus
	eventBus, err := events.NewNATSEventBus(cfg.NATS.URL)
	if err != nil {
		logger.Fatal("Failed to connect to NATS", zap.Error(err))
	}
	defer eventBus.Close()

	// Setup event subscriptions
	setupEventSubscriptions(eventBus)

	logger.Info("Worker is ready and listening for events")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Worker shutting down...")
	logger.Info("Worker shutdown complete")
}

// setupEventSubscriptions sets up NATS event subscriptions for background processing
func setupEventSubscriptions(eventBus events.EventBus) {
	// Subscribe to UserRegistered events for background processing
	err := eventBus.Subscribe(domain.UserRegisteredEvent, handleUserRegisteredBackground)
	if err != nil {
		logger.Error("Failed to subscribe to UserRegistered events", zap.Error(err))
	}

	// Subscribe to UserUpgradedToSeller events
	err = eventBus.Subscribe(domain.UserUpgradedToSellerEvent, handleUserUpgradedToSeller)
	if err != nil {
		logger.Error("Failed to subscribe to UserUpgradedToSeller events", zap.Error(err))
	}

	logger.Info("Worker event subscriptions setup complete")
}

// Background event handlers
func handleUserRegisteredBackground(ctx context.Context, event *events.Event) error {
	logger.Info("Worker handling UserRegistered event",
		zap.String("event_id", event.ID),
		zap.String("user_id", event.AggregateID))

	var userData domain.UserRegistered
	if err := events.ParseEventData(event, &userData); err != nil {
		return err
	}

	// Background processing tasks:
	// 1. Send verification email
	// 2. Add to marketing automation
	// 3. Update analytics
	// 4. Initialize user preferences
	// 5. Create user folder structure

	// Simulate some background work
	time.Sleep(100 * time.Millisecond)

	logger.Info("Worker completed UserRegistered background processing",
		zap.String("user_email", userData.Email))

	return nil
}

func handleUserUpgradedToSeller(ctx context.Context, event *events.Event) error {
	logger.Info("Worker handling UserUpgradedToSeller event",
		zap.String("event_id", event.ID),
		zap.String("user_id", event.AggregateID))

	var userData domain.UserUpgradedToSeller
	if err := events.ParseEventData(event, &userData); err != nil {
		return err
	}

	// Background processing for new sellers:
	// 1. Send seller onboarding emails
	// 2. Create seller analytics dashboard
	// 3. Initialize seller tools and resources
	// 4. Notify admin for verification review
	// 5. Add to seller notification channels

	logger.Info("Worker completed UserUpgradedToSeller background processing",
		zap.String("user_email", userData.Email),
		zap.String("business_name", userData.BusinessName))

	return nil
}
