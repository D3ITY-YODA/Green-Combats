package broker

import (
	"context"
	"log/slog"
	"sync"
)

// Event is the interface all broker events must implement.
type Event interface {
	Topic() string
}

// Handler processes a single event.
type Handler func(ctx context.Context, evt Event) error

// Broker is an in-process publish/subscribe event bus.
type Broker struct {
	mu     sync.RWMutex
	subs   map[string][]handlerEntry
	logger *slog.Logger
}

type handlerEntry struct {
	id      int
	handler Handler
}

var nextID int

// New creates a new Broker.
func New(logger *slog.Logger) *Broker {
	return &Broker{
		subs:   make(map[string][]handlerEntry),
		logger: logger,
	}
}

// Subscribe registers a handler for events matching the given topic.
// Use topic "*" to receive all events.
// Returns an unsubscribe function.
func (b *Broker) Subscribe(topic string, h Handler) func() {
	b.mu.Lock()
	defer b.mu.Unlock()

	nextID++
	id := nextID
	b.subs[topic] = append(b.subs[topic], handlerEntry{id: id, handler: h})

	return func() {
		b.unsubscribe(topic, id)
	}
}

func (b *Broker) unsubscribe(topic string, id int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	entries := b.subs[topic]
	for i, e := range entries {
		if e.id == id {
			b.subs[topic] = append(entries[:i], entries[i+1:]...)
			return
		}
	}
}

// Publish sends an event to all subscribers of the matching topic.
// Handlers for the specific topic and for "*" are both invoked.
// Errors from handlers are logged but do not stop delivery to other handlers.
func (b *Broker) Publish(ctx context.Context, evt Event) {
	topic := evt.Topic()

	b.mu.RLock()
	// Copy the slice to avoid holding the lock during handler execution
	specific := make([]handlerEntry, len(b.subs[topic]))
	copy(specific, b.subs[topic])
	wildcard := make([]handlerEntry, len(b.subs["*"]))
	copy(wildcard, b.subs["*"])
	b.mu.RUnlock()

	all := append(specific, wildcard...)
	for _, entry := range all {
		if err := entry.handler(ctx, evt); err != nil {
			b.logger.Error("broker: handler error",
				"topic", topic,
				"handler_id", entry.id,
				"err", err,
			)
		}
	}
}

// PublishAsync sends an event in a background goroutine.
func (b *Broker) PublishAsync(ctx context.Context, evt Event) {
	go b.Publish(ctx, evt)
}

// Topics returns the list of topics with at least one subscriber.
func (b *Broker) Topics() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	topics := make([]string, 0, len(b.subs))
	for t, entries := range b.subs {
		if len(entries) > 0 {
			topics = append(topics, t)
		}
	}
	return topics
}
