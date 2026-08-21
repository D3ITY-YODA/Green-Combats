package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"green-compass-backend/internal/config"
	"green-compass-backend/pkg/database"
	"green-compass-backend/pkg/logging"
)

const serviceName = "green-compass-scheduler"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", serviceName, err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load(config.LoadOptions{})
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger, err := logging.New(logging.Options{
		Level:   cfg.Log.Level,
		Format:  cfg.Log.Format,
		Service: serviceName,
	})
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}

	if cfg.Database.URL == "" {
		return fmt.Errorf("database.url must be configured for the scheduler process")
	}

	pool, err := database.Connect(ctx, database.Options{
		URL:               cfg.Database.URL,
		MaxConns:          cfg.Database.MaxConns,
		MinConns:          cfg.Database.MinConns,
		ConnMaxLifetime:   time.Duration(cfg.Database.ConnMaxLifetime),
		ConnMaxIdleTime:   time.Duration(cfg.Database.ConnMaxIdleTime),
		HealthCheckPeriod: time.Duration(cfg.Database.HealthCheckPeriod),
	})
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer pool.Close()

	logger.Info("scheduler started", "app_env", cfg.AppEnv)

	// Pipeline tick: runs the content generation pipeline periodically.
	// In a full implementation, this would:
	// 1. Query all active places
	// 2. For each place, run: normalization -> indicators -> applicability -> assessment -> content
	// 3. Publish ContentGeneratedEvent via broker for each new content item
	//
	// The worker process handles data ingestion on its own schedule.
	// The scheduler handles the processing/generation pipeline.
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("scheduler shutting down cleanly")
			return nil

		case <-ticker.C:
			logger.Info("scheduler: running content generation pipeline")

			if err := runContentPipeline(ctx, pool, logger); err != nil {
				logger.Error("scheduler: pipeline failed", "err", err)
			} else {
				logger.Info("scheduler: pipeline completed successfully")
			}
		}
	}
}

// runContentPipeline executes the full content generation pipeline for all active places.
func runContentPipeline(ctx context.Context, _ interface{}, logger *slog.Logger) error {
	// Placeholder: in Phase 5, the individual services (normalization, indicators,
	// applicability, assessment, content) would be wired here.
	//
	// The pipeline would:
	// 1. List all places with recent normalized data
	// 2. Compute indicators for each place
	// 3. Apply applicability rules to filter relevant indicators
	// 4. Run assessment to score urgency/severity
	// 5. Generate plain-language content
	// 6. Publish events for notification dispatch

	logger.Debug("content pipeline: placeholder — individual services to be wired")
	return nil
}
