// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package events

import (
	"log"
	"time"
)

// EventType represents the type of event
type EventType string

const (
	EventTypePush     EventType = "push"
	EventTypePull     EventType = "pull"
	EventTypePromote  EventType = "promote"
	EventTypeRollback EventType = "rollback"
	EventTypeDelete   EventType = "delete"
)

// Event represents an event in the system
type Event struct {
	ID        int                    `json:"id"`
	Type      EventType              `json:"type"`
	Project   string                 `json:"project"`
	App       string                 `json:"app"`
	Version   string                 `json:"version,omitempty"`
	AgentID   string                 `json:"agent_id,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// EventBus handles event publishing and subscription
type EventBus interface {
	Publish(event *Event) error
	Subscribe(eventType EventType, handler EventHandler) error
}

// EventHandler handles events
type EventHandler func(event *Event) error

// MemoryEventBus is a simple in-memory event bus implementation
type MemoryEventBus struct {
	handlers map[EventType][]EventHandler
}

// NewMemoryEventBus creates a new memory event bus
func NewMemoryEventBus() *MemoryEventBus {
	return &MemoryEventBus{
		handlers: make(map[EventType][]EventHandler),
	}
}

// Publish publishes an event
func (b *MemoryEventBus) Publish(event *Event) error {
	handlers := b.handlers[event.Type]
	for _, handler := range handlers {
		if err := handler(event); err != nil {
			// Log error but continue with other handlers
			// Use standard library log for now (can be replaced with structured logger if needed)
			log.Printf("Event handler error for event type %s (project=%s, app=%s, version=%s): %v",
				event.Type, event.Project, event.App, event.Version, err)
		}
	}
	return nil
}

// Subscribe subscribes to an event type
func (b *MemoryEventBus) Subscribe(eventType EventType, handler EventHandler) error {
	b.handlers[eventType] = append(b.handlers[eventType], handler)
	return nil
}
