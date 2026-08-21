package health_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"green-compass-backend/internal/health"
	"green-compass-backend/pkg/clock"
)

func newTestRouter(t *testing.T, clk clock.Clock) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	handler := health.NewHandler(health.NewService("test-version", clk))
	handler.RegisterRoutes(router)
	return router
}

func TestHandler_HealthReturnsOK(t *testing.T) {
	at := time.Date(2026, 8, 21, 8, 0, 0, 0, time.UTC)
	router := newTestRouter(t, clock.NewFixed(at))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Errorf("content-type = %q, want application/json", ct)
	}

	var envelope struct {
		Data health.StatusResponse `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode response body: %v\nbody: %s", err, rec.Body.String())
	}
	body := envelope.Data
	if body.Status != health.StatusOK {
		t.Errorf("status = %q, want %q", body.Status, health.StatusOK)
	}
	if body.Version != "test-version" {
		t.Errorf("version = %q, want test-version", body.Version)
	}
	if body.UptimeSeconds != 0 {
		t.Errorf("uptime_seconds = %d, want 0 with fixed clock", body.UptimeSeconds)
	}
}

func TestHandler_UptimeAdvancesWithClock(t *testing.T) {
	at := time.Date(2026, 8, 21, 8, 0, 0, 0, time.UTC)
	clk := clock.NewFixed(at)
	router := newTestRouter(t, clk)

	clk.Advance(2 * time.Minute)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	var envelope struct {
		Data health.StatusResponse `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode response body: %v", err)
	}
	if envelope.Data.UptimeSeconds != 120 {
		t.Errorf("uptime_seconds = %d, want 120", envelope.Data.UptimeSeconds)
	}
}

func TestService_StatusAcceptsContext(t *testing.T) {
	svc := health.NewService("v1", clock.NewFixed(time.Time{}))

	resp := svc.Status(context.Background())

	if resp.Status != health.StatusOK {
		t.Errorf("status = %q, want %q", resp.Status, health.StatusOK)
	}
}
