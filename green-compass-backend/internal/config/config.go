package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/goccy/go-yaml"

	"green-compass-backend/pkg/logging"
)

const (
	EnvLocal      = "local"
	EnvDev        = "dev"
	EnvStaging    = "staging"
	EnvProduction = "production"
)

type Duration time.Duration

func (d *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	parsed, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", s, err)
	}
	if parsed <= 0 {
		return fmt.Errorf("invalid duration %q: must be positive", s)
	}
	*d = Duration(parsed)
	return nil
}

type Server struct {
	Host            string   `yaml:"host"`
	Port            int      `yaml:"port"`
	ReadTimeout     Duration `yaml:"read_timeout"`
	WriteTimeout    Duration `yaml:"write_timeout"`
	IdleTimeout     Duration `yaml:"idle_timeout"`
	ShutdownTimeout Duration `yaml:"shutdown_timeout"`
}

type Logging struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type Database struct {
	URL               string   `yaml:"url"`
	MaxConns          int32    `yaml:"max_conns"`
	MinConns          int32    `yaml:"min_conns"`
	ConnMaxLifetime   Duration `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime   Duration `yaml:"conn_max_idle_time"`
	HealthCheckPeriod Duration `yaml:"health_check_period"`
}

type Config struct {
	AppEnv   string   `yaml:"app_env"`
	Server   Server   `yaml:"server"`
	Log      Logging  `yaml:"logging"`
	Database Database `yaml:"database"`
}

type LoadOptions struct {
	Path   string
	Lookup func(string) (string, bool)
}

func Load(opts LoadOptions) (*Config, error) {
	lookup := opts.Lookup
	if lookup == nil {
		lookup = os.LookupEnv
	}

	cfg := Defaults()

	data, err := readFile(opts.Path, lookup)
	if err != nil {
		return nil, err
	}
	if data != nil {
		if err := yaml.UnmarshalWithOptions(data, cfg, yaml.Strict()); err != nil {
			return nil, fmt.Errorf("parse config: %w", err)
		}
	}

	if err := applyEnvOverrides(cfg, lookup); err != nil {
		return nil, err
	}

	if errs := cfg.Validate(); len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return cfg, nil
}

func Defaults() *Config {
	return &Config{
		AppEnv: EnvLocal,
		Server: Server{
			Host:            "",
			Port:            8080,
			ReadTimeout:     Duration(10 * time.Second),
			WriteTimeout:    Duration(10 * time.Second),
			IdleTimeout:     Duration(120 * time.Second),
			ShutdownTimeout: Duration(15 * time.Second),
		},
		Log: Logging{
			Level:  "info",
			Format: "json",
		},
		Database: Database{
			URL:               "",
			MaxConns:          10,
			MinConns:          0,
			ConnMaxLifetime:   Duration(30 * time.Minute),
			ConnMaxIdleTime:   Duration(5 * time.Minute),
			HealthCheckPeriod: Duration(1 * time.Minute),
		},
	}
}

func (c *Config) Validate() []error {
	var errs []error

	switch c.AppEnv {
	case EnvLocal, EnvDev, EnvStaging, EnvProduction:
	default:
		errs = append(errs, fmt.Errorf("app_env: invalid value %q: must be one of local, dev, staging, production", c.AppEnv))
	}

	if c.Server.Port < 1 || c.Server.Port > 65535 {
		errs = append(errs, fmt.Errorf("server.port: %d is outside valid range 1-65535", c.Server.Port))
	}

	for name, d := range map[string]Duration{
		"server.read_timeout":     c.Server.ReadTimeout,
		"server.write_timeout":    c.Server.WriteTimeout,
		"server.idle_timeout":     c.Server.IdleTimeout,
		"server.shutdown_timeout": c.Server.ShutdownTimeout,
	} {
		if d <= 0 {
			errs = append(errs, fmt.Errorf("%s: must be positive, got %s", name, time.Duration(d)))
		}
	}

	if _, err := logging.ParseLevel(c.Log.Level); err != nil {
		errs = append(errs, fmt.Errorf("logging.level: %w", err))
	}
	if _, err := logging.ParseFormat(c.Log.Format); err != nil {
		errs = append(errs, fmt.Errorf("logging.format: %w", err))
	}

	if c.Database.MaxConns < 0 {
		errs = append(errs, fmt.Errorf("database.max_conns: must not be negative, got %d", c.Database.MaxConns))
	}
	if c.Database.MinConns < 0 {
		errs = append(errs, fmt.Errorf("database.min_conns: must not be negative, got %d", c.Database.MinConns))
	}
	if c.Database.MaxConns > 0 && c.Database.MinConns > c.Database.MaxConns {
		errs = append(errs, fmt.Errorf("database.min_conns: %d exceeds max_conns %d", c.Database.MinConns, c.Database.MaxConns))
	}

	return errs
}

func readFile(path string, lookup func(string) (string, bool)) ([]byte, error) {
	explicit := false
	if path == "" {
		path, _ = lookup("CONFIG_FILE")
	}
	if path != "" {
		explicit = true
	} else {
		env, _ := lookup("GC_APP_ENV")
		if env == "" {
			env = EnvLocal
		}
		path = filepath.Join("configs", "config."+env+".yaml")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) && !explicit {
			return nil, nil
		}
		return nil, fmt.Errorf("read config file %q: %w", path, err)
	}
	return data, nil
}

func applyEnvOverrides(cfg *Config, lookup func(string) (string, bool)) error {
	var errs []error

	applyString(&cfg.AppEnv, lookup, "GC_APP_ENV")
	applyString(&cfg.Server.Host, lookup, "GC_SERVER_HOST")

	if v, ok := lookup("GC_SERVER_PORT"); ok && v != "" {
		port, err := strconv.Atoi(v)
		if err != nil {
			errs = append(errs, fmt.Errorf("env GC_SERVER_PORT: invalid integer %q", v))
		} else {
			cfg.Server.Port = port
		}
	}

	for key, dst := range map[string]*Duration{
		"GC_SERVER_READ_TIMEOUT":     &cfg.Server.ReadTimeout,
		"GC_SERVER_WRITE_TIMEOUT":    &cfg.Server.WriteTimeout,
		"GC_SERVER_IDLE_TIMEOUT":     &cfg.Server.IdleTimeout,
		"GC_SERVER_SHUTDOWN_TIMEOUT": &cfg.Server.ShutdownTimeout,
	} {
		v, ok := lookup(key)
		if !ok || v == "" {
			continue
		}
		d, err := time.ParseDuration(v)
		if err != nil {
			errs = append(errs, fmt.Errorf("env %s: invalid duration %q", key, v))
			continue
		}
		*dst = Duration(d)
	}

	applyString(&cfg.Log.Level, lookup, "GC_LOG_LEVEL")
	applyString(&cfg.Log.Format, lookup, "GC_LOG_FORMAT")

	applyString(&cfg.Database.URL, lookup, "GC_DATABASE_URL")

	if v, ok := lookup("GC_DATABASE_MAX_CONNS"); ok && v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			errs = append(errs, fmt.Errorf("env GC_DATABASE_MAX_CONNS: invalid integer %q", v))
		} else {
			cfg.Database.MaxConns = int32(n)
		}
	}
	if v, ok := lookup("GC_DATABASE_MIN_CONNS"); ok && v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			errs = append(errs, fmt.Errorf("env GC_DATABASE_MIN_CONNS: invalid integer %q", v))
		} else {
			cfg.Database.MinConns = int32(n)
		}
	}

	for key, dst := range map[string]*Duration{
		"GC_DATABASE_CONN_MAX_LIFETIME": &cfg.Database.ConnMaxLifetime,
		"GC_DATABASE_IDLE_TIME":         &cfg.Database.ConnMaxIdleTime,
		"GC_DATABASE_HEALTH_CHECK":      &cfg.Database.HealthCheckPeriod,
	} {
		v, ok := lookup(key)
		if !ok || v == "" {
			continue
		}
		d, err := time.ParseDuration(v)
		if err != nil || d <= 0 {
			errs = append(errs, fmt.Errorf("env %s: invalid duration %q", key, v))
			continue
		}
		*dst = Duration(d)
	}

	return errors.Join(errs...)
}

func applyString(dst *string, lookup func(string) (string, bool), key string) {
	if v, ok := lookup(key); ok && v != "" {
		*dst = v
	}
}
