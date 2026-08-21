package clock_test

import (
	"testing"
	"time"

	"green-compass-backend/pkg/clock"
)

func TestNew_ReturnsSystemTime(t *testing.T) {
	before := time.Now()
	got := clock.New().Now()
	after := time.Now()

	if got.Before(before) || got.After(after) {
		t.Fatalf("Now() = %v, want a time between %v and %v", got, before, after)
	}
}

func TestFixed_NowIsStable(t *testing.T) {
	at := time.Date(2026, 8, 21, 9, 0, 0, 0, time.UTC)
	fixed := clock.NewFixed(at)

	for i := 0; i < 3; i++ {
		if got := fixed.Now(); !got.Equal(at) {
			t.Fatalf("Now() = %v, want %v (call %d)", got, at, i+1)
		}
	}
}

func TestFixed_Advance(t *testing.T) {
	start := time.Date(2026, 8, 21, 9, 0, 0, 0, time.UTC)
	fixed := clock.NewFixed(start)

	fixed.Advance(90 * time.Second)

	want := start.Add(90 * time.Second)
	if got := fixed.Now(); !got.Equal(want) {
		t.Fatalf("Now() = %v, want %v after Advance(90s)", got, want)
	}
}
