// Package queue provides a RabbitMQ consumer with auto-reconnect and
// routing-key-based handler dispatch. It is designed for consuming internal
// event messages such as privilege changes, login events, etc.
//
// Lifecycle:
//
//	consumer, _ := queue.NewConsumer(&cfg.RabbitMQ)
//	consumer.RegisterHandler("privilege.role.updated", onRoleUpdated)
//	consumer.Start(ctx)
//	defer consumer.Stop()
package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"

	"github.com/phoenix-agent-go/internal/config"
)

const (
	defaultReconnectDelay       = 3 * time.Second
	defaultMaxReconnectAttempts = 5
	defaultPrefetchCount        = 10
	consumerQueueName           = "phoenix.consumer"
)

// MessageHandler is the signature for a message handler function.
// It receives the enriched context and the raw message body.
type MessageHandler func(ctx context.Context, body []byte) error

// EventMessage is the common event envelope consumed from RabbitMQ.
// Type is the routing key / event type; Payload is the event-specific JSON.
type EventMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// CacheInvalidator is a callback for clearing cached data when events demand it.
// The consumer does not directly depend on the cache layer — the caller wires
// this during construction.
type CacheInvalidator func(ctx context.Context, eventType string, payload []byte)

// Consumer is a RabbitMQ consumer that dispatches messages to registered
// handlers based on the routing key. It supports auto-reconnect and
// graceful shutdown.
type Consumer struct {
	cfg    *config.RabbitMQConfig
	logger *zap.Logger

	conn    *amqp.Connection
	channel *amqp.Channel

	handlers map[string]MessageHandler
	mu       sync.RWMutex

	invalidator CacheInvalidator

	connNotify   chan *amqp.Error
	done         chan struct{}
	reconnecting bool
	reconnectMu  sync.Mutex
}

// NewConsumer creates a new Consumer, establishes the RabbitMQ connection,
// and declares the topic exchange. It does NOT start consuming — call Start()
// to begin processing messages.
func NewConsumer(cfg *config.RabbitMQConfig) (*Consumer, error) {
	c := &Consumer{
		cfg:      cfg,
		logger:   zap.L().Named("queue.consumer"),
		handlers: make(map[string]MessageHandler),
		done:     make(chan struct{}),
	}

	// Apply defaults for zero-value config fields.
	if cfg.ReconnectDelay == 0 {
		cfg.ReconnectDelay = defaultReconnectDelay
	}
	if cfg.MaxReconnectAttempts == 0 {
		cfg.MaxReconnectAttempts = defaultMaxReconnectAttempts
	}
	if cfg.PrefetchCount == 0 {
		cfg.PrefetchCount = defaultPrefetchCount
	}

	if err := c.connect(); err != nil {
		return nil, fmt.Errorf("consumer initial connect: %w", err)
	}

	// Pre-register privilege event handlers.
	c.registerDefaultHandlers()

	return c, nil
}

// SetCacheInvalidator attaches a cache invalidation callback. The consumer
// calls this for events that require cache clearing (e.g., role updates).
func (c *Consumer) SetCacheInvalidator(fn CacheInvalidator) {
	c.invalidator = fn
}

// ──────────────────────────── Connection Management ────────────────────────────

func (c *Consumer) connect() error {
	c.logger.Info("connecting to RabbitMQ", zap.String("addr", c.cfg.Addr))

	conn, err := amqp.Dial(c.cfg.Addr)
	if err != nil {
		return fmt.Errorf("amqp dial: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("amqp channel: %w", err)
	}

	// Declare the topic exchange (idempotent).
	if err := channel.ExchangeDeclare(
		c.cfg.Exchange, // name
		"topic",        // type
		true,           // durable
		false,          // auto-delete
		false,          // internal
		false,          // no-wait
		nil,            // args
	); err != nil {
		channel.Close()
		conn.Close()
		return fmt.Errorf("exchange declare: %w", err)
	}

	// Declare the consumer queue.
	_, err = channel.QueueDeclare(
		consumerQueueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return fmt.Errorf("queue declare: %w", err)
	}

	if err := channel.Qos(c.cfg.PrefetchCount, 0, false); err != nil {
		channel.Close()
		conn.Close()
		return fmt.Errorf("set qos: %w", err)
	}

	c.conn = conn
	c.channel = channel
	c.connNotify = conn.NotifyClose(make(chan *amqp.Error, 1))

	c.logger.Info("RabbitMQ connected", zap.String("exchange", c.cfg.Exchange))
	return nil
}

// ensureBindings binds all registered routing keys to the consumer queue.
// Caller must hold c.mu (at least read lock).
func (c *Consumer) ensureBindings() error {
	for routingKey := range c.handlers {
		if err := c.channel.QueueBind(
			consumerQueueName, // queue
			routingKey,        // routing key
			c.cfg.Exchange,    // exchange
			false,             // no-wait
			nil,               // args
		); err != nil {
			return fmt.Errorf("bind %q: %w", routingKey, err)
		}
	}
	return nil
}

func (c *Consumer) close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil && !c.conn.IsClosed() {
		c.conn.Close()
	}
}

// ──────────────────────────── Handler Registration ────────────────────────────

// RegisterHandler associates a routing key with a handler function.
// If a handler already exists for the key, it is replaced.
//
// Handlers should be registered before Start() is called to ensure bindings
// are in place. Registration after Start() is safe but requires a manual
// re-bind (the existing consumer won't pick up new routing keys until
// the next reconnect).
func (c *Consumer) RegisterHandler(routingKey string, handler MessageHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[routingKey] = handler

	// If already connected, bind immediately.
	if c.channel != nil {
		if err := c.channel.QueueBind(
			consumerQueueName,
			routingKey,
			c.cfg.Exchange,
			false,
			nil,
		); err != nil {
			c.logger.Warn("failed to bind routing key", zap.String("key", routingKey), zap.Error(err))
		}
	}
}

// ──────────────────────────── Start / Stop ────────────────────────────

// Start begins consuming messages. It blocks until the consumer is shut down
// (via Stop or context cancellation). Run this in a goroutine if needed.
func (c *Consumer) Start(ctx context.Context) error {
	// Bind all registered routing keys.
	c.mu.RLock()
	if err := c.ensureBindings(); err != nil {
		c.mu.RUnlock()
		return fmt.Errorf("ensure bindings: %w", err)
	}
	c.mu.RUnlock()

	deliveries, err := c.channel.Consume(
		consumerQueueName,
		"",    // consumer tag (auto-generated)
		false, // auto-ack: we ack manually after handler succeeds
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("start consume: %w", err)
	}

	c.logger.Info("consumer started")

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("consumer context cancelled, stopping")
			return nil
		case <-c.done:
			c.logger.Info("consumer stopped via Stop()")
			return nil
		case amqpErr := <-c.connNotify:
			c.logger.Warn("connection lost, attempting reconnect", zap.Error(amqpErr))
			c.reconnect(ctx, deliveries)
			// After reconnect, get a new deliveries channel.
			var consumeErr error
			deliveries, consumeErr = c.channel.Consume(
				consumerQueueName, "", false, false, false, false, nil,
			)
			if consumeErr != nil {
				return fmt.Errorf("re-consume after reconnect: %w", consumeErr)
			}
		case msg, ok := <-deliveries:
			if !ok {
				c.logger.Warn("delivery channel closed, attempting reconnect")
				c.reconnect(ctx, deliveries)
				var consumeErr error
				deliveries, consumeErr = c.channel.Consume(
					consumerQueueName, "", false, false, false, false, nil,
				)
				if consumeErr != nil {
					return fmt.Errorf("re-consume after channel close: %w", consumeErr)
				}
				continue
			}
			c.handleDelivery(msg)
		}
	}
}

// Stop initiates a graceful shutdown of the consumer.
func (c *Consumer) Stop() error {
	c.logger.Info("stopping consumer")
	select {
	case <-c.done:
		// Already stopped.
	default:
		close(c.done)
	}
	c.close()
	return nil
}

// ──────────────────────────── Message Processing ────────────────────────────

func (c *Consumer) handleDelivery(msg amqp.Delivery) {
	eventType := msg.RoutingKey

	c.mu.RLock()
	handler, ok := c.handlers[eventType]
	c.mu.RUnlock()

	if !ok {
		c.logger.Warn("no handler for routing key, nacking",
			zap.String("routingKey", eventType))
		msg.Nack(false, false)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := handler(ctx, msg.Body); err != nil {
		c.logger.Error("handler failed, nacking (requeue)",
			zap.String("event", eventType),
			zap.Error(err),
		)
		msg.Nack(false, true) // requeue for retry
		return
	}

	msg.Ack(false)
}

// ──────────────────────────── Auto-Reconnect ────────────────────────────

func (c *Consumer) reconnect(ctx context.Context, oldDeliveries <-chan amqp.Delivery) {
	c.reconnectMu.Lock()
	defer c.reconnectMu.Unlock()

	if !c.reconnecting {
		c.reconnecting = true
		defer func() { c.reconnecting = false }()
	} else {
		return // Already reconnecting elsewhere.
	}

	c.close()

	for attempt := 1; attempt <= c.cfg.MaxReconnectAttempts; attempt++ {
		c.logger.Info("reconnecting", zap.Int("attempt", attempt),
			zap.Duration("delay", c.cfg.ReconnectDelay))

		select {
		case <-ctx.Done():
			return
		case <-c.done:
			return
		case <-time.After(c.cfg.ReconnectDelay):
		}

		if err := c.connect(); err != nil {
			c.logger.Warn("reconnect failed", zap.Int("attempt", attempt), zap.Error(err))
			continue
		}

		// Re-bind all routing keys.
		c.mu.RLock()
		if err := c.ensureBindings(); err != nil {
			c.mu.RUnlock()
			c.logger.Warn("re-bind failed", zap.Error(err))
			c.close()
			continue
		}
		c.mu.RUnlock()

		// Drain remaining messages from the old delivery channel.
		go drainDeliveries(oldDeliveries)

		c.logger.Info("reconnected successfully", zap.Int("attempt", attempt))
		return
	}

	c.logger.Error("max reconnect attempts reached, giving up")
}

func drainDeliveries(ch <-chan amqp.Delivery) {
	for range ch {
		// Drain until closed.
	}
}

// ──────────────────────────── Default Handlers ────────────────────────────

func (c *Consumer) registerDefaultHandlers() {
	c.RegisterHandler("privilege.user.created", c.handleUserCreated)
	c.RegisterHandler("privilege.role.updated", c.handleRoleUpdated)
	c.RegisterHandler("privilege.login.success", c.handleLoginSuccess)
}

func (c *Consumer) handleUserCreated(ctx context.Context, body []byte) error {
	c.logger.Info("privilege.user.created event received")

	var msg EventMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		c.logger.Warn("failed to unmarshal user.created event", zap.Error(err))
		return nil // Do not requeue malformed messages.
	}

	c.logger.Info("user created", zap.String("payload", string(msg.Payload)))
	return nil
}

func (c *Consumer) handleRoleUpdated(ctx context.Context, body []byte) error {
	c.logger.Info("privilege.role.updated event received")

	var msg EventMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		c.logger.Warn("failed to unmarshal role.updated event", zap.Error(err))
		return nil
	}

	c.logger.Info("role updated, invalidating permission cache",
		zap.String("payload", string(msg.Payload)))

	// Invalidate permission cache if a cache invalidator is registered.
	if c.invalidator != nil {
		c.invalidator(ctx, "privilege.role.updated", msg.Payload)
	}

	return nil
}

func (c *Consumer) handleLoginSuccess(ctx context.Context, body []byte) error {
	c.logger.Info("privilege.login.success event received")

	var msg EventMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		c.logger.Warn("failed to unmarshal login.success event", zap.Error(err))
		return nil
	}

	c.logger.Info("login success", zap.String("payload", string(msg.Payload)))
	return nil
}
