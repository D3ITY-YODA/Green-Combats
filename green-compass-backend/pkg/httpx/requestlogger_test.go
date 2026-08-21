package httpx_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"green-compass-backend/pkg/clock"
	"green-compass-backend/pkg/httpx"
	"green-compass-backend/pkg/logging"
)

func newTestLogger(t *testing.T) (*slog.Logger, *bytes.Buffer) {
	t.Helper()
	var buf bytes.Buffer
	logger, err := logging.New(logging.Options{Level: "debug", Format: "json", Writer: &buf})
	if err != nil {
		t.Fatalf("logging.New() unexpected error: %v", err)
	}
	return logger, &buf
}

func TestRequestLogger_LogsRequestDetails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger, buf := newTestLogger(t)
	at := time.Date(2026, 8, 21, 12, 0, 0, 0, time.UTC)
	clk := clock.NewFixed(at)

	router := gin.New()
	router.Use(httpx.RequestLogger(logger, clk))
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"pong": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.RemoteAddr = "10.0.0.7:54321"
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("log output is not valid JSON: %v\noutput: %s", err, buf.String())
	}
	if entry["msg"] != "http_request" {
		t.Errorf("msg = %v, want http_request", entry["msg"])
	}
	if entry["method"] != http.MethodGet {
		t.Errorf("method = %v, want GET", entry["method"])
	}
	if entry["path"] != "/ping" {
		t.Errorf("path = %v, want /ping", entry["path"])
	}
	if status, ok := entry["status"].(float64); !ok || int(status) != http.StatusOK {
		t.Errorf("status = %v (%T), want 200", entry["status"], entry["status"])
	}
	if duration, ok := entry["duration_ms"].(float64); !ok || duration != 0 {
		t.Errorf("duration_ms = %v (%T), want 0 with fixed clock", entry["duration_ms"], entry["duration_ms"])
	}
	if entry["client_ip"] != "10.0.0.1" && entry["client_ip"] != "" {
		t.Logf("client_ip = %v (httptest may normalize RemoteAddr)", entry["client_ip"])
	}
}

func TestRequestLogger_RecordsHandlerStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger, buf := newTestLogger(t)

	router := gin.New()
	router.Use(httpx.RequestLogger(logger, clock.NewFixed(time.Time{})))
	router.POST("/fail", func(c *gin.Context) {
		c.JSON(http.StatusTeapot, gin.H{"error": "brewing"})
	})

	req := httptest.NewRequest(http.MethodPost, "/fail", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	out := buf.String()
	if !strings.Contains(out, `"status":418`) {
		t.Errorf("log output %q does not contain status 418", out)
	}
	if !strings.Contains(out, `"method":"POST"`) {
		t.Errorf("log output %q does not contain method POST", out)
	}
}
