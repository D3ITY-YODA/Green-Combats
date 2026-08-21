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
	if cfg.Database.URL != "" {
		t.Errorf("Database.URL = %q, want empty default", cfg.Database.URL)
	}
	if cfg.Database.MaxConns != 10 {
		t.Errorf("Database.MaxConns = %d, want 10", cfg.Database.MaxConns)
	}
	if time.Duration(cfg.Database.ConnMaxLifetime) != 30*time.Minute {
		t.Errorf("Database.ConnMaxLifetime = %s, want 30m", time.Duration(cfg.Database.ConnMaxLifetime))
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
database:
  url: postgres://user:pass@db:5432/app
  max_conns: 25
  min_conns: 2
  conn_max_lifetime: 45m
  conn_max_idle_time: 8m
  health_check_period: 2m
auth:
  secret: "staging-secret-0123456789abcdef0123456789abcdef"
`)

	cfg, err := config.Load(config.LoadOptions{Path: path, Lookup: lookupFrom(nil)})
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
	if cfg.Database.URL != "postgres://user:pass@db:5432/app" {
		t.Errorf("Database.URL = %q, want value from file", cfg.Database.URL)
	}
	if cfg.Database.MaxConns != 25 || cfg.Database.MinConns != 2 {
		t.Errorf("Database pool sizes = %d/%d, want 25/2", cfg.Database.MaxConns, cfg.Database.MinConns)
	}
	if time.Duration(cfg.Database.ConnMaxIdleTime) != 8*time.Minute {
		t.Errorf("Database.ConnMaxIdleTime = %s, want 8m", time.Duration(cfg.Database.ConnMaxIdleTime))
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
		"GC_SERVER_PORT":                "9100",
		"GC_LOG_LEVEL":                  "debug",
		"GC_APP_ENV":                    "dev",
		"GC_DATABASE_URL":               "postgres://env:env@db:5432/env",
		"GC_DATABASE_MAX_CONNS":         "30",
		"GC_DATABASE_CONN_MAX_LIFETIME": "1h",
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
	if cfg.Database.URL != "postgres://env:env@db:5432/env" {
		t.Errorf("Database.URL = %q, want env override", cfg.Database.URL)
	}
	if cfg.Database.MaxConns != 30 {
		t.Errorf("Database.MaxConns = %d, want env override 30", cfg.Database.MaxConns)
	}
	if time.Duration(cfg.Database.ConnMaxLifetime) != time.Hour {
		t.Errorf("Database.ConnMaxLifetime = %s, want env override 1h", time.Duration(cfg.Database.ConnMaxLifetime))
	}
}

func TestLoad_AuthSection(t *testing.T) {
	dir := t.TempDir()
	path := writeConfig(t, dir, "auth.yaml", `
app_env: local
auth:
  secret: "0123456789abcdef0123456789abcdef"
  access_token_ttl: 10m
  refresh_token_ttl: 48h
  issuer: custom-issuer
`)

	cfg, err := config.Load(config.LoadOptions{Path: path, Lookup: lookupFrom(nil)})
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}
	if cfg.Auth.Secret != "0123456789abcdef0123456789abcdef" {
		t.Errorf("Auth.Secret = %q, want value from file", cfg.Auth.Secret)
	}
	if time.Duration(cfg.Auth.AccessTokenTTL) != 10*time.Minute {
		t.Errorf("Auth.AccessTokenTTL = %s, want 10m", time.Duration(cfg.Auth.AccessTokenTTL))
	}
	if time.Duration(cfg.Auth.RefreshTokenTTL) != 48*time.Hour {
		t.Errorf("Auth.RefreshTokenTTL = %s, want 48h", time.Duration(cfg.Auth.RefreshTokenTTL))
	}
	if cfg.Auth.Issuer != "custom-issuer" {
		t.Errorf("Auth.Issuer = %q, want custom-issuer", cfg.Auth.Issuer)
	}
}

func TestLoad_AuthEnvOverrides(t *testing.T) {
	dir := t.TempDir()
	path := writeConfig(t, dir, "base.yaml", "app_env: local\n")

	env := map[string]string{
		"GC_AUTH_SECRET":      "env-secret-0123456789abcdef0123456789abcdef",
		"GC_AUTH_ACCESS_TTL":  "5m",
		"GC_AUTH_REFRESH_TTL": "24h",
		"GC_AUTH_ISSUER":      "env-issuer",
	}

	cfg, err := config.Load(config.LoadOptions{Path: path, Lookup: lookupFrom(env)})
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}
	if cfg.Auth.Secret != env["GC_AUTH_SECRET"] {
		t.Errorf("Auth.Secret = %q, want env override", cfg.Auth.Secret)
	}
	if time.Duration(cfg.Auth.AccessTokenTTL) != 5*time.Minute {
		t.Errorf("Auth.AccessTokenTTL = %s, want env override 5m", time.Duration(cfg.Auth.AccessTokenTTL))
	}
	if time.Duration(cfg.Auth.RefreshTokenTTL) != 24*time.Hour {
		t.Errorf("Auth.RefreshTokenTTL = %s, want env override 24h", time.Duration(cfg.Auth.RefreshTokenTTL))
	}
	if cfg.Auth.Issuer != "env-issuer" {
		t.Errorf("Auth.Issuer = %q, want env override env-issuer", cfg.Auth.Issuer)
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

	_, err := config.Load(config.LoadOptions{Path: path, Lookup: lookupFrom(nil)})
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
		{
			name:    "bad env database duration",
			yaml:    "",
			env:     map[string]string{"GC_DATABASE_IDLE_TIME": "soon"},
			wantErr: "GC_DATABASE_IDLE_TIME",
		},
		{
			name:    "negative env max conns",
			yaml:    "",
			env:     map[string]string{"GC_DATABASE_MAX_CONNS": "-3"},
			wantErr: "database.max_conns",
		},
		{
			name:    "min conns above max",
			yaml:    "database:\n  max_conns: 5\n  min_conns: 9\n",
			wantErr: "database.min_conns",
		},
		{
			name:    "negative min conns in file",
			yaml:    "database:\n  min_conns: -1\n",
			wantErr: "database.min_conns",
		},
		{
			name:    "empty auth secret in staging",
			yaml:    "app_env: staging\nauth:\n  secret: \"\"\n",
			wantErr: "auth.secret",
		},
		{
			name:    "short auth secret in production",
			yaml:    "app_env: production\nauth:\n  secret: too-short\n",
			wantErr: "auth.secret",
		},
		{
			name:    "refresh ttl not above access ttl",
			yaml:    "auth:\n  access_token_ttl: 30m\n  refresh_token_ttl: 15m\n",
			wantErr: "auth.refresh_token_ttl",
		},
		{
			name:    "empty issuer",
			yaml:    "auth:\n  issuer: \"\"\n",
			wantErr: "auth.issuer",
		},
		{
			name:    "bad env auth duration",
			yaml:    "",
			env:     map[string]string{"GC_AUTH_ACCESS_TTL": "soon"},
			wantErr: "GC_AUTH_ACCESS_TTL",
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
		Database: config.Database{
			MaxConns: 4,
			MinConns: 9,
		},
	}

	errs := cfg.Validate()
	if len(errs) < 5 {
		t.Fatalf("Validate() returned %d errors, want at least 5: %v", len(errs), errs)
	}
}
