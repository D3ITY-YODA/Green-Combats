package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"green-compass-backend/internal/config"
)

func lookupFrom(m map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		v, ok := m[key]
		return v, ok
	}
}

func writeConfig(t *testing.T, dir, name, contents string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("write config file: %v", err)
	}
	return path
}

func TestDefaults_AreValid(t *testing.T) {
	cfg := config.Defaults()
	if errs := cfg.Validate(); len(errs) != 0 {
		t.Fatalf("defaults failed validation: %v", errs)
	}

	if cfg.AppEnv != config.EnvLocal {
		t.Errorf("AppEnv = %q, want %q", cfg.AppEnv, config.EnvLocal)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("Port = %d, want 8080", cfg.Server.Port)
	}
	if time.Duration(cfg.Server.ReadTimeout) != 10*time.Second {
		t.Errorf("ReadTimeout = %s, want 10s", time.Duration(cfg.Server.ReadTimeout))
	}
	if cfg.Log.Level != "info" || cfg.Log.Format != "json" {
		t.Errorf("Log = %+v, want level info format json", cfg.Log)
	}
}

func TestLoad_FromFile(t *testing.T) {
	dir := t.TempDir()
	path := writeConfig(t, dir, "custom.yaml", `
app_env: staging
server:
  host: 127.0.0.1
  port: 9000
  read_timeout: 5s
  write_timeout: 6s
  idle_timeout: 60s
  shutdown_timeout: 20s
logging:
  level: debug
  format: text
`)

	cfg, err := config.Load(config.LoadOptions{Path: path})
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if cfg.AppEnv != config.EnvStaging {
		t.Errorf("AppEnv = %q, want %q", cfg.AppEnv, config.EnvStaging)
	}
	if cfg.Server.Port != 9000 {
		t.Errorf("Port = %d, want 9000", cfg.Server.Port)
	}
	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("Host = %q, want 127.0.0.1", cfg.Server.Host)
	}
	if time.Duration(cfg.Server.IdleTimeout) != time.Minute {
		t.Errorf("IdleTimeout = %s, want 1m", time.Duration(cfg.Server.IdleTimeout))
	}
	if cfg.Log.Level != "debug" || cfg.Log.Format != "text" {
		t.Errorf("Log = %+v, want level debug format text", cfg.Log)
	}
}

func TestLoad_EnvOverridesFile(t *testing.T) {
	dir := t.TempDir()
	path := writeConfig(t, dir, "base.yaml", `
app_env: local
server:
  port: 9000
logging:
  level: info
  format: json
`)

	env := map[string]string{
		"GC_SERVER_PORT": "9100",
		"GC_LOG_LEVEL":   "debug",
		"GC_APP_ENV":     "dev",
	}

	cfg, err := config.Load(config.LoadOptions{Path: path, Lookup: lookupFrom(env)})
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if cfg.Server.Port != 9100 {
		t.Errorf("Port = %d, want env override 9100", cfg.Server.Port)
	}
	if cfg.Log.Level != "debug" {
		t.Errorf("Log.Level = %q, want env override debug", cfg.Log.Level)
	}
	if cfg.AppEnv != config.EnvDev {
		t.Errorf("AppEnv = %q, want env override dev", cfg.AppEnv)
	}
}

func TestLoad_UnsetEnvDoesNotOverride(t *testing.T) {
	dir := t.TempDir()
	path := writeConfig(t, dir, "base.yaml", `
server:
  port: 9000
`)

	env := map[string]string{
		"GC_SERVER_PORT": "",
		"GC_LOG_FORMAT":  "",
	}

	cfg, err := config.Load(config.LoadOptions{Path: path, Lookup: lookupFrom(env)})
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if cfg.Server.Port != 9000 {
		t.Errorf("Port = %d, want 9000 from file (empty env must not override)", cfg.Server.Port)
	}
	if cfg.Log.Format != "json" {
		t.Errorf("Log.Format = %q, want default json", cfg.Log.Format)
	}
}

func TestLoad_MissingConventionalFile_FallsBackToDefaults(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	cfg, err := config.Load(config.LoadOptions{Lookup: lookupFrom(nil)})
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("Port = %d, want default 8080", cfg.Server.Port)
	}
}

func TestLoad_DiscoveredByAppEnv(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	if err := os.MkdirAll(filepath.Join(dir, "configs"), 0o755); err != nil {
		t.Fatalf("mkdir configs: %v", err)
	}
	writeConfig(t, filepath.Join(dir, "configs"), "config.dev.yaml", `
app_env: dev
server:
  port: 9200
`)

	cfg, err := config.Load(config.LoadOptions{
		Lookup: lookupFrom(map[string]string{"GC_APP_ENV": "dev"}),
	})
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}
	if cfg.Server.Port != 9200 {
		t.Errorf("Port = %d, want 9200 from configs/config.dev.yaml", cfg.Server.Port)
	}
}

func TestLoad_ExplicitMissingFile_IsError(t *testing.T) {
	_, err := config.Load(config.LoadOptions{Path: "/nonexistent/config.yaml"})
	if err == nil {
		t.Fatal("Load() expected error for explicitly missing file, got nil")
	}
	if !strings.Contains(err.Error(), "/nonexistent/config.yaml") {
		t.Errorf("error = %v, want it to mention the missing path", err)
	}
}

func TestLoad_StrictDecoding_RejectsUnknownKeys(t *testing.T) {
	dir := t.TempDir()
	path := writeConfig(t, dir, "typo.yaml", `
server:
  prot: 9000
`)

	_, err := config.Load(config.LoadOptions{Path: path})
	if err == nil {
		t.Fatal("Load() expected error for unknown key, got nil")
	}
}

func TestLoad_Errors(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		env     map[string]string
		wantErr string
	}{
		{
			name:    "invalid duration in file",
			yaml:    "server:\n  read_timeout: fast\n",
			wantErr: "invalid duration",
		},
		{
			name:    "zero duration in file",
			yaml:    "server:\n  read_timeout: 0s\n",
			wantErr: "must be positive",
		},
		{
			name:    "bad app_env",
			yaml:    "app_env: prod\n",
			wantErr: "app_env",
		},
		{
			name:    "port below range",
			yaml:    "server:\n  port: 0\n",
			wantErr: "server.port",
		},
		{
			name:    "port above range",
			yaml:    "server:\n  port: 70000\n",
			wantErr: "server.port",
		},
		{
			name:    "bad log level",
			yaml:    "logging:\n  level: loud\n",
			wantErr: "logging.level",
		},
		{
			name:    "bad log format",
			yaml:    "logging:\n  format: xml\n",
			wantErr: "logging.format",
		},
		{
			name:    "bad env port",
			yaml:    "",
			env:     map[string]string{"GC_SERVER_PORT": "abc"},
			wantErr: "GC_SERVER_PORT",
		},
		{
			name:    "bad env duration",
			yaml:    "",
			env:     map[string]string{"GC_SERVER_IDLE_TIMEOUT": "soon"},
			wantErr: "GC_SERVER_IDLE_TIMEOUT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			opts := config.LoadOptions{Lookup: lookupFrom(tt.env)}
			if tt.yaml != "" {
				opts.Path = writeConfig(t, dir, "cfg.yaml", tt.yaml)
			} else {
				t.Chdir(dir)
			}

			_, err := config.Load(opts)
			if err == nil {
				t.Fatal("Load() expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error = %v, want it to contain %q", err, tt.wantErr)
			}
		})
	}
}

func TestValidate_AggregatesMultipleErrors(t *testing.T) {
	cfg := &config.Config{
		AppEnv: "nope",
		Server: config.Server{Port: -1},
		Log:    config.Logging{Level: "loud", Format: "xml"},
	}

	errs := cfg.Validate()
	if len(errs) < 4 {
		t.Fatalf("Validate() returned %d errors, want at least 4: %v", len(errs), errs)
	}
}
