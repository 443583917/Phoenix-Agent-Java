package event

import (
	"context"
	"encoding/json"
	"sync"

	"go.uber.org/zap"
)

type Handler func(ctx context.Context, event *Event) error

type EventBus struct {
	handlers map[string][]Handler
	mu       sync.RWMutex
	logger   *zap.Logger
	publishFn func(ctx context.Context, event *Event) error
}

func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]Handler),
		logger:   zap.L().Named("event.bus"),
	}
}

func (b *EventBus) SetPublisher(fn func(ctx context.Context, event *Event) error) {
	b.publishFn = fn
}

func (b *EventBus) Subscribe(eventType string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

func (b *EventBus) Publish(ctx context.Context, event *Event) error {
	if b.publishFn != nil {
		return b.publishFn(ctx, event)
	}
	return b.dispatch(ctx, event)
}

func (b *EventBus) dispatch(ctx context.Context, event *Event) error {
	b.mu.RLock()
	handlers := b.handlers[event.Type]
	b.mu.RUnlock()

	for _, h := range handlers {
		if err := h(ctx, event); err != nil {
			b.logger.Error("event handler failed",
				zap.String("type", event.Type),
				zap.Error(err),
			)
		}
	}
	return nil
}

func (b *EventBus) HandleMessage(ctx context.Context, eventType string, payload json.RawMessage) error {
	event := &Event{
		Type:    eventType,
		Payload: payload,
	}
	return b.dispatch(ctx, event)
}
