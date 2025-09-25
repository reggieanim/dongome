package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"dongome/internal/users/domain"
	"dongome/pkg/events"
)

// This example demonstrates how UserRegistered events flow through NATS
// between different bounded contexts in a real-world scenario.

func main() {
	// Initialize NATS event bus
	eventBus, err := events.NewNATSEventBus("nats://localhost:4222")
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer eventBus.Close()

	fmt.Println("🚀 Starting NATS Event Flow Demonstration")
	fmt.Println("📡 Connected to NATS at nats://localhost:4222")

	// Setup event subscribers (simulating different bounded contexts)
	setupEventSubscribers(eventBus)

	// Wait a moment for subscriptions to be ready
	time.Sleep(1 * time.Second)

	// Simulate a user registration event
	simulateUserRegistration(eventBus)

	// Keep the program running to see event processing
	fmt.Println("\n⏳ Waiting for events to be processed...")
	time.Sleep(5 * time.Second)

	fmt.Println("\n✅ Event flow demonstration completed!")
}

func setupEventSubscribers(eventBus events.EventBus) {
	fmt.Println("\n📋 Setting up event subscribers...")

	// 1. Notifications bounded context subscriber
	err := eventBus.Subscribe(domain.UserRegisteredEvent, handleUserRegisteredForNotifications)
	if err != nil {
		log.Printf("Failed to subscribe notifications handler: %v", err)
	} else {
		fmt.Println("   ✓ Notifications handler subscribed")
	}

	// 2. Analytics bounded context subscriber
	err = eventBus.Subscribe(domain.UserRegisteredEvent, handleUserRegisteredForAnalytics)
	if err != nil {
		log.Printf("Failed to subscribe analytics handler: %v", err)
	} else {
		fmt.Println("   ✓ Analytics handler subscribed")
	}

	// 3. Marketing bounded context subscriber
	err = eventBus.Subscribe(domain.UserRegisteredEvent, handleUserRegisteredForMarketing)
	if err != nil {
		log.Printf("Failed to subscribe marketing handler: %v", err)
	} else {
		fmt.Println("   ✓ Marketing handler subscribed")
	}
}

func simulateUserRegistration(eventBus events.EventBus) {
	fmt.Println("\n👤 Simulating user registration...")

	// Create a sample UserRegistered event
	userRegisteredData := domain.UserRegistered{
		UserID:            "user-123e4567-e89b-12d3-a456-426614174000",
		Email:             "john.doe@example.com",
		FirstName:         "John",
		LastName:          "Doe",
		Role:              "buyer",
		VerificationToken: "verification-token-abc123",
		Timestamp:         time.Now(),
	}

	// Create the event
	event, err := events.NewEvent(
		domain.UserRegisteredEvent,
		userRegisteredData.UserID,
		userRegisteredData,
	)
	if err != nil {
		log.Fatalf("Failed to create event: %v", err)
	}

	// Publish the event to NATS
	ctx := context.Background()
	if err := eventBus.Publish(ctx, event); err != nil {
		log.Fatalf("Failed to publish event: %v", err)
	}

	fmt.Printf("   ✓ UserRegistered event published for user: %s %s (%s)\n",
		userRegisteredData.FirstName,
		userRegisteredData.LastName,
		userRegisteredData.Email)
}

// Event Handlers for different bounded contexts

// handleUserRegisteredForNotifications simulates the Notifications bounded context
func handleUserRegisteredForNotifications(ctx context.Context, event *events.Event) error {
	fmt.Println("\n📧 [NOTIFICATIONS CONTEXT] Processing UserRegistered event")

	var userData domain.UserRegistered
	if err := events.ParseEventData(event, &userData); err != nil {
		return fmt.Errorf("failed to parse event data: %w", err)
	}

	// Simulate notification processing
	fmt.Printf("   📨 Preparing welcome email for: %s\n", userData.Email)
	fmt.Printf("   📱 Queuing SMS verification for user: %s\n", userData.UserID)
	fmt.Printf("   🔔 Setting up push notification preferences\n")

	// Simulate some processing time
	time.Sleep(500 * time.Millisecond)

	fmt.Printf("   ✅ Notifications processed for user: %s\n", userData.FirstName)
	return nil
}

// handleUserRegisteredForAnalytics simulates the Analytics bounded context
func handleUserRegisteredForAnalytics(ctx context.Context, event *events.Event) error {
	fmt.Println("\n📊 [ANALYTICS CONTEXT] Processing UserRegistered event")

	var userData domain.UserRegistered
	if err := events.ParseEventData(event, &userData); err != nil {
		return fmt.Errorf("failed to parse event data: %w", err)
	}

	// Simulate analytics processing
	fmt.Printf("   📈 Recording user registration metric\n")
	fmt.Printf("   🌍 Analyzing user demographics\n")
	fmt.Printf("   📋 Creating user journey tracking\n")
	fmt.Printf("   🎯 Initializing conversion funnel for user: %s\n", userData.UserID)

	// Simulate some processing time
	time.Sleep(300 * time.Millisecond)

	fmt.Printf("   ✅ Analytics data recorded for user: %s\n", userData.Email)
	return nil
}

// handleUserRegisteredForMarketing simulates the Marketing bounded context
func handleUserRegisteredForMarketing(ctx context.Context, event *events.Event) error {
	fmt.Println("\n🎯 [MARKETING CONTEXT] Processing UserRegistered event")

	var userData domain.UserRegistered
	if err := events.ParseEventData(event, &userData); err != nil {
		return fmt.Errorf("failed to parse event data: %w", err)
	}

	// Simulate marketing processing
	fmt.Printf("   📝 Adding user to onboarding campaign\n")
	fmt.Printf("   🎁 Preparing welcome bonus for: %s\n", userData.FirstName)
	fmt.Printf("   📊 Segmenting user profile\n")
	fmt.Printf("   📧 Scheduling marketing emails\n")

	// Simulate some processing time
	time.Sleep(400 * time.Millisecond)

	fmt.Printf("   ✅ Marketing setup completed for user: %s\n", userData.Email)
	return nil
}

// Helper function to pretty print event data
func prettyPrintEvent(event *events.Event) {
	eventJSON, _ := json.MarshalIndent(event, "", "  ")
	fmt.Printf("Event Details:\n%s\n", eventJSON)
}
