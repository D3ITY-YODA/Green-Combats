package logging_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"green-compass-backend/pkg/logging"
)

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "empty defaults to info", input: ""},
		{name: "debug", input: "debug"},
		{name: "info uppercase", input: "INFO"},
		{name: "warn", input: "warn"},
		{name: "warning alias", input: "warning"},
		{name: "error mixed case", input: "Error"},
		{name: "padded whitespace", input: "  debug  "},
		{name: "unknown level", input: "verbose", wantErr: true},
		{name: "numeric level rejected", input: "4", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, err := logging.ParseLevel(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParseLevel(%q) expected error, got nil", tt.input)
				}
				if !strings.Contains(err.Error(), "invalid log level") {
					t.Fatalf("ParseLevel(%q) error = %v, want it to mention invalid log level", tt.input, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseLevel(%q) unexpected error: %v", tt.input, err)
			}
			if tt.input == "" && level != 0 {
				t.Fatalf("ParseLevel(\"\") = %v, want default info level", level)
			}
		})
	}
}

func TestParseFormat(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    logging.Format
		wantErr bool
	}{
		{name: "empty defaults to text", input: "", want: logging.FormatText},
		{name: "text", input: "text", want: logging.FormatText},
		{name: "json uppercase", input: "JSON", want: logging.FormatJSON},
		{name: "unknown format", input: "xml", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := logging.ParseFormat(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParseFormat(%q) expected error, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseFormat(%q) unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("ParseFormat(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNew_JSONOutput(t *testing.T) {
	var buf bytes.Buffer
	logger, err := logging.New(logging.Options{
		Level:   "info",
		Format:  "json",
		Service: "test-service",
		Writer:  &buf,
	})
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	logger.Info("request handled", "path", "/health", "status", 200)

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("log output is not valid JSON: %v\noutput: %s", err, buf.String())
	}
	if entry["msg"] != "request handled" {
		t.Errorf("msg = %v, want %q", entry["msg"], "request handled")
	}
	if entry["service"] != "test-service" {
		t.Errorf("service = %v, want %q", entry["service"], "test-service")
	}
	if entry["path"] != "/health" {
		t.Errorf("path = %v, want %q", entry["path"], "/health")
	}
	if status, ok := entry["status"].(float64); !ok || status != 200 {
		t.Errorf("status = %v (%T), want 200", entry["status"], entry["status"])
	}
}

func TestNew_TextOutput(t *testing.T) {
	var buf bytes.Buffer
	logger, err := logging.New(logging.Options{Format: "text", Writer: &buf})
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	logger.Warn("cache_miss", "key", "places:1")

	out := buf.String()
	if !strings.Contains(out, "msg=cache_miss") {
		t.Errorf("text output %q does not contain msg=cache_miss", out)
	}
	if !strings.Contains(out, "key=places:1") {
		t.Errorf("text output %q does not contain key=places:1", out)
	}
	if strings.HasPrefix(out, "{") {
		t.Errorf("text output %q looks like JSON", out)
	}
}

func TestNew_LevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger, err := logging.New(logging.Options{Level: "warn", Format: "json", Writer: &buf})
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	logger.Info("should not appear")
	if buf.Len() != 0 {
		t.Errorf("info record was written despite warn level: %s", buf.String())
	}

	logger.Error("should appear")
	if !strings.Contains(buf.String(), "should appear") {
		t.Errorf("error record missing from output: %s", buf.String())
	}
}

func TestNew_InvalidOptions(t *testing.T) {
	tests := []struct {
		name string
		opts logging.Options
	}{
		{name: "bad level", opts: logging.Options{Level: "loud"}},
		{name: "bad format", opts: logging.Options{Format: "csv"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := logging.New(tt.opts); err == nil {
				t.Fatalf("New(%+v) expected error, got nil", tt.opts)
			}
		})
	}
}
