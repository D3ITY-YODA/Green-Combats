package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"green-compass-backend/internal/config"
	"green-compass-backend/internal/connectors"
	"green-compass-backend/internal/ingestion"
	"green-compass-backend/internal/places"
	"green-compass-backend/internal/sources"
	"green-compass-backend/pkg/database"
	"green-compass-backend/pkg/logging"
)

const serviceName = "green-compass-worker"

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
		return fmt.Errorf("database.url must be configured for the worker process")
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

	openMeteo := connectors.NewOpenMeteo(connectors.HTTPOptions{})
	nasaPower := connectors.NewNASAPower(connectors.HTTPOptions{})
	registry, err := connectors.NewRegistry(openMeteo, nasaPower)
	if err != nil {
		return fmt.Errorf("init connector registry: %w", err)
	}

	ingestionRepo := ingestion.NewRepository(pool)
	ingestionSvc := ingestion.NewService(ingestionRepo, registry)
	sourcesSvc := sources.NewService(sources.NewRepository(pool))
	placesRepo := places.NewRepository(pool)

	logger.Info("ingestion worker started", "app_env", cfg.AppEnv)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("ingestion worker shutting down cleanly")
			return nil

		case <-ticker.C:
			enabledSources, err := sourcesSvc.Enabled(ctx)
			if err != nil {
				logger.Error("failed to load enabled data sources", "err", err)
				continue
			}

			for _, src := range enabledSources {
				if ctx.Err() != nil {
					break
				}
				// Process places for each enabled source
				now := time.Now().UTC()
				from := now.Add(-src.PollInterval)
				to := now

				// Pull for saved/configured places (using default limit for discovery)
				allPlaces, err := placesRepo.Search(ctx, "", 100)
				if err != nil {
					logger.Error("failed to list places for ingestion", "source", src.Code, "err", err)
					continue
				}

				for _, p := range allPlaces {
					if ctx.Err() != nil {
						break
					}
					result, err := ingestionSvc.Ingest(ctx, ingestion.IngestRequest{
						Source: src,
						Place:  p,
						From:   from,
						To:     to,
					})
					if err != nil {
						logger.Warn("ingestion run failed", "source", src.Code, "place_id", p.ID, "err", err)
					} else {
						logger.Info("ingestion run completed",
							"source", src.Code,
							"place_id", p.ID,
							"fetched", result.FetchedCount,
							"landed", result.LandedCount,
						)
					}
				}
			}
		}
	}
}
