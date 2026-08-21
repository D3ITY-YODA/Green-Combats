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
	"green-compass-backend/internal/notifications"
	"green-compass-backend/pkg/broker"
	"green-compass-backend/pkg/database"
	"green-compass-backend/pkg/logging"
)

const serviceName = "green-compass-notifier"

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
		return fmt.Errorf("database.url must be configured for the notifier process")
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

	notifRepo := notifications.NewRepository(pool)
	notifSvc := notifications.NewService(notifRepo, logger)

	b := broker.New(logger)

	// Subscribe to content.generated events
	b.Subscribe(broker.TopicContentGenerated, func(ctx context.Context, evt broker.Event) error {
		contentEvt, ok := evt.(broker.ContentGeneratedEvent)
		if !ok {
			return fmt.Errorf("unexpected event type: %T", evt)
		}

		logger.Info("content generated notification triggered",
			"content_id", contentEvt.ContentID,
			"place_id", contentEvt.PlaceID,
			"place_name", contentEvt.PlaceName,
		)

		// In a real implementation, this would:
		// 1. Query users who have saved this place
		// 2. Check their notification preferences
		// 3. Create notification records for each eligible user+channel
		// 4. Dispatch via the appropriate adapter (push/SMS/USSD)
		logger.Info("notifier: would dispatch notification for content update",
			"place", contentEvt.PlaceName,
			"language", contentEvt.Language,
		)

		return nil
	})

	// Subscribe to observation.submitted events
	b.Subscribe(broker.TopicObservationSubmit, func(ctx context.Context, evt broker.Event) error {
		obsEvt, ok := evt.(broker.ObservationSubmittedEvent)
		if !ok {
			return fmt.Errorf("unexpected event type: %T", evt)
		}

		logger.Info("observation submitted notification triggered",
			"observation_id", obsEvt.ObservationID,
			"user_id", obsEvt.UserID,
			"place_id", obsEvt.PlaceID,
		)
		return nil
	})

	// Subscribe to report.verified events
	b.Subscribe(broker.TopicReportVerified, func(ctx context.Context, evt broker.Event) error {
		reportEvt, ok := evt.(broker.ReportVerifiedEvent)
		if !ok {
			return fmt.Errorf("unexpected event type: %T", evt)
		}

		logger.Info("report verified notification triggered",
			"observation_id", reportEvt.ObservationID,
			"org_id", reportEvt.OrgID,
			"verified_by", reportEvt.VerifiedBy,
		)
		return nil
	})

	// Start the pending notification dispatcher
	go dispatchPending(ctx, notifSvc, logger)

	logger.Info("notifier started", "app_env", cfg.AppEnv)
	<-ctx.Done()
	logger.Info("notifier shutting down cleanly")
	return nil
}

// dispatchPending periodically picks up pending notifications and attempts delivery.
func dispatchPending(ctx context.Context, svc *notifications.Service, logger *slog.Logger) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// In a real implementation, this would:
			// 1. Fetch pending notifications from the DB
			// 2. Dispatch each via the appropriate adapter
			// 3. Mark as sent or failed
			logger.Debug("dispatcher tick: checking for pending notifications")
		}
	}
}
