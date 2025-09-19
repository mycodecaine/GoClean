package events

import (
	"context"
	"goclean/internal/domain/entities"
)

// DomainEventPublisher represents a domain event publisher interface
type DomainEventPublisher interface {
	Publish(ctx context.Context, events []entities.DomainEvent) error
}

// DomainEventHandler represents a domain event handler interface
type DomainEventHandler interface {
	Handle(ctx context.Context, event entities.DomainEvent) error
	CanHandle(event entities.DomainEvent) bool
}

// DomainEventDispatcher handles dispatching domain events to appropriate handlers
type DomainEventDispatcher struct {
	handlers  []DomainEventHandler
	publisher DomainEventPublisher
}

// NewDomainEventDispatcher creates a new domain event dispatcher
func NewDomainEventDispatcher(publisher DomainEventPublisher) *DomainEventDispatcher {
	return &DomainEventDispatcher{
		handlers:  make([]DomainEventHandler, 0),
		publisher: publisher,
	}
}

// RegisterHandler registers a domain event handler
func (d *DomainEventDispatcher) RegisterHandler(handler DomainEventHandler) {
	d.handlers = append(d.handlers, handler)
}

// DispatchEvents dispatches domain events from an aggregate root
func (d *DomainEventDispatcher) DispatchEvents(ctx context.Context, aggregateRoot *entities.AggregateRoot) error {
	events := aggregateRoot.DomainEvents()
	if len(events) == 0 {
		return nil
	}

	// Handle events locally first
	for _, event := range events {
		for _, handler := range d.handlers {
			if handler.CanHandle(event) {
				if err := handler.Handle(ctx, event); err != nil {
					return err
				}
			}
		}
	}

	// Publish events for external consumption
	if d.publisher != nil {
		if err := d.publisher.Publish(ctx, events); err != nil {
			return err
		}
	}

	// Clear events after successful dispatch
	aggregateRoot.ClearDomainEvents()
	return nil
}
