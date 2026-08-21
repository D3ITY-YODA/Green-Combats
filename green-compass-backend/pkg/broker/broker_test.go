package broker

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type testEvent struct {
	topic string
	data  string
}

func (e testEvent) Topic() string { return e.topic }

func newTestBroker() *Broker {
	return New(slog.Default())
}

func TestPublishSubscribe(t *testing.T) {
	b := newTestBroker()
	var received []Event

	unsub := b.Subscribe("test.topic", func(ctx context.Context, evt Event) error {
		received = append(received, evt)
		return nil
	})
	defer unsub()

	b.Publish(context.Background(), testEvent{topic: "test.topic", data: "hello"})
	b.Publish(context.Background(), testEvent{topic: "test.topic", data: "world"})

	if len(received) != 2 {
		t.Fatalf("expected 2 events, got %d", len(received))
	}
	if received[0].(testEvent).data != "hello" {
		t.Errorf("first event: got %q, want %q", received[0].(testEvent).data, "hello")
	}
	if received[1].(testEvent).data != "world" {
		t.Errorf("second event: got %q, want %q", received[1].(testEvent).data, "world")
	}
}

func TestUnsubscribe(t *testing.T) {
	b := newTestBroker()
	var count atomic.Int32

	unsub := b.Subscribe("t", func(ctx context.Context, evt Event) error {
		count.Add(1)
		return nil
	})

	b.Publish(context.Background(), testEvent{topic: "t"})
	if count.Load() != 1 {
		t.Fatalf("expected 1, got %d", count.Load())
	}

	unsub()
	b.Publish(context.Background(), testEvent{topic: "t"})
	if count.Load() != 1 {
		t.Fatalf("expected still 1 after unsubscribe, got %d", count.Load())
	}
}

func TestWildcardSubscriber(t *testing.T) {
	b := newTestBroker()
	var topics []string
	var mu sync.Mutex

	b.Subscribe("*", func(ctx context.Context, evt Event) error {
		mu.Lock()
		topics = append(topics, evt.Topic())
		mu.Unlock()
		return nil
	})

	b.Publish(context.Background(), testEvent{topic: "a"})
	b.Publish(context.Background(), testEvent{topic: "b"})

	mu.Lock()
	defer mu.Unlock()
	if len(topics) != 2 {
		t.Fatalf("expected 2 events, got %d", len(topics))
	}
}

func TestHandlerErrorLoggedNotFatal(t *testing.T) {
	b := newTestBroker()
	var count atomic.Int32

	b.Subscribe("t", func(ctx context.Context, evt Event) error {
		return errors.New("boom")
	})
	b.Subscribe("t", func(ctx context.Context, evt Event) error {
		count.Add(1)
		return nil
	})

	b.Publish(context.Background(), testEvent{topic: "t"})
	if count.Load() != 1 {
		t.Fatalf("second handler should still run, got count %d", count.Load())
	}
}

func TestPublishAsync(t *testing.T) {
	b := newTestBroker()
	var received atomic.Bool

	b.Subscribe("t", func(ctx context.Context, evt Event) error {
		received.Store(true)
		return nil
	})

	b.PublishAsync(context.Background(), testEvent{topic: "t"})

	// Give goroutine time to execute
	time.Sleep(50 * time.Millisecond)
	if !received.Load() {
		t.Fatal("async handler was not called")
	}
}

func TestTopics(t *testing.T) {
	b := newTestBroker()

	unsub1 := b.Subscribe("alpha", func(ctx context.Context, evt Event) error { return nil })
	unsub2 := b.Subscribe("beta", func(ctx context.Context, evt Event) error { return nil })
	defer unsub1()
	defer unsub2()

	topics := b.Topics()
	if len(topics) != 2 {
		t.Fatalf("expected 2 topics, got %d: %v", len(topics), topics)
	}
}

func TestNoSubscribersDoesNotPanic(t *testing.T) {
	b := newTestBroker()
	// Should not panic
	b.Publish(context.Background(), testEvent{topic: "nonexistent"})
}

func TestMultipleSubscribersSameTopic(t *testing.T) {
	b := newTestBroker()
	var count atomic.Int32

	for i := 0; i < 5; i++ {
		b.Subscribe("t", func(ctx context.Context, evt Event) error {
			count.Add(1)
			return nil
		})
	}

	b.Publish(context.Background(), testEvent{topic: "t"})
	if count.Load() != 5 {
		t.Fatalf("expected 5 handlers called, got %d", count.Load())
	}
}
