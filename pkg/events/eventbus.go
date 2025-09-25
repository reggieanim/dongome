package events

import (
	"context"
	"encoding/json"
	"time"

	"dongome/pkg/logger"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// Event represents a domain event
type Event struct {
	ID          string            `json:"id"`
	Type        string            `json:"type"`
	AggregateID string            `json:"aggregate_id"`
	Data        json.RawMessage   `json:"data"`
	Metadata    map[string]string `json:"metadata"`
	Timestamp   time.Time         `json:"timestamp"`
}

// EventBus defines the interface for event publishing and subscribing
type EventBus interface {
	Publish(ctx context.Context, event *Event) error
	Subscribe(eventType string, handler EventHandler) error
	Close() error
}

// EventHandler defines the signature for event handlers
type EventHandler func(ctx context.Context, event *Event) error

// NATSEventBus implements EventBus using NATS
type NATSEventBus struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

// NewNATSEventBus creates a new NATS event bus
func NewNATSEventBus(url string) (*NATSEventBus, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	// Create JetStream context
	js, err := conn.JetStream()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Create stream if it doesn't exist
	streamName := "DOMAIN_EVENTS"
	_, err = js.StreamInfo(streamName)
	if err == nats.ErrStreamNotFound {
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{"events.>"},
			Storage:  nats.FileStorage,
			Replicas: 1,
		})
		if err != nil {
			conn.Close()
			return nil, err
		}
	} else if err != nil {
		conn.Close()
		return nil, err
	}

	return &NATSEventBus{
		conn: conn,
		js:   js,
	}, nil
}

// Publish publishes an event to NATS
func (eb *NATSEventBus) Publish(ctx context.Context, event *Event) error {
	// Set event metadata
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Serialize event
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// Publish to NATS subject
	subject := "events." + event.Type
	_, err = eb.js.PublishAsync(subject, data)
	if err != nil {
		logger.Error("Failed to publish event",
			zap.String("event_id", event.ID),
			zap.String("event_type", event.Type),
			zap.Error(err))
		return err
	}

	logger.Info("Event published",
		zap.String("event_id", event.ID),
		zap.String("event_type", event.Type),
		zap.String("aggregate_id", event.AggregateID))

	return nil
}

// Subscribe subscribes to events of a specific type
func (eb *NATSEventBus) Subscribe(eventType string, handler EventHandler) error {
	subject := "events." + eventType

	_, err := eb.js.Subscribe(subject, func(msg *nats.Msg) {
		// Parse event
		var event Event
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			logger.Error("Failed to unmarshal event",
				zap.String("subject", msg.Subject),
				zap.Error(err))
			msg.Nak()
			return
		}

		// Handle event with timeout context
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := handler(ctx, &event); err != nil {
			logger.Error("Event handler failed",
				zap.String("event_id", event.ID),
				zap.String("event_type", event.Type),
				zap.Error(err))
			msg.Nak()
			return
		}

		logger.Info("Event handled successfully",
			zap.String("event_id", event.ID),
			zap.String("event_type", event.Type))

		msg.Ack()
	}, nats.Durable("dongome-"+eventType))

	if err != nil {
		return err
	}

	logger.Info("Subscribed to event type", zap.String("event_type", eventType))
	return nil
}

// Close closes the NATS connection
func (eb *NATSEventBus) Close() error {
	if eb.conn != nil {
		eb.conn.Close()
	}
	return nil
}

// NewEvent creates a new event
func NewEvent(eventType string, aggregateID string, data interface{}) (*Event, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &Event{
		ID:          uuid.New().String(),
		Type:        eventType,
		AggregateID: aggregateID,
		Data:        dataBytes,
		Metadata:    make(map[string]string),
		Timestamp:   time.Now(),
	}, nil
}

// ParseEventData parses event data into the provided type
func ParseEventData(event *Event, target interface{}) error {
	return json.Unmarshal(event.Data, target)
}
