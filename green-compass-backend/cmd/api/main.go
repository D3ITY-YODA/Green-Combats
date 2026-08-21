package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"green-compass-backend/internal/audit"
	gcctx "green-compass-backend/internal/context"
	"green-compass-backend/internal/auth"
	"green-compass-backend/internal/config"
	"green-compass-backend/internal/console"
	"green-compass-backend/internal/health"
	"green-compass-backend/internal/integrations"
	"green-compass-backend/internal/notifications"
	"green-compass-backend/internal/observations"
	"green-compass-backend/internal/places"
	"green-compass-backend/internal/preferences"
	"green-compass-backend/internal/projects"
	"green-compass-backend/internal/reporting"
	"green-compass-backend/internal/reports"
	"green-compass-backend/internal/sources"
	"green-compass-backend/internal/updates"
	"green-compass-backend/internal/users"
	"green-compass-backend/pkg/clock"
	"green-compass-backend/pkg/database"
	"green-compass-backend/pkg/httpx"
	"green-compass-backend/pkg/logging"
	"green-compass-backend/pkg/metrics"
	"green-compass-backend/pkg/storage"
)

const serviceName = "green-compass-api"

var version = "dev"

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

	if cfg.AppEnv == config.EnvProduction || cfg.AppEnv == config.EnvStaging {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	if cfg.Database.URL == "" {
		return fmt.Errorf("database.url must be configured for the api process")
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

	userSvc := users.NewService(users.NewRepository(pool))
	authRepo := auth.NewRepository(pool)
	authSvc, err := auth.NewService(userSvc, authRepo, clock.New(), auth.Options{
		Secret:     cfg.Auth.Secret,
		Issuer:     cfg.Auth.Issuer,
		AccessTTL:  time.Duration(cfg.Auth.AccessTokenTTL),
		RefreshTTL: time.Duration(cfg.Auth.RefreshTokenTTL),
	})
	if err != nil {
		return fmt.Errorf("init auth service: %w", err)
	}

	router := gin.New()
	registry := metrics.NewRegistry()
	router.Use(
		gin.Recovery(),
		httpx.RequestIDMiddleware(),
		httpx.RequestLogger(logger, clock.New()),
		metrics.Middleware(registry),
		httpx.IdempotencyMiddleware(),
		httpx.NewRateLimiter(120, 240).Middleware(), // 120 req/min per client, burst 240
	)
	router.GET("/metrics", metrics.Handler(registry))
	router.Use(httpx.CORSMiddleware(httpx.CORSOptions{AllowedOrigins: cfg.CORS.AllowedOrigins}))

	healthService := health.NewService(version, clock.New())
	health.NewHandler(healthService).RegisterRoutes(router)
	auth.NewHandler(authSvc).RegisterRoutes(router)

	placeSvc := places.NewService(places.NewRepository(pool))
	places.NewHandler(placeSvc, authSvc).RegisterRoutes(router)

	preferencesSvc := preferences.NewService(preferences.NewRepository(pool))
	preferences.NewHandler(preferencesSvc).RegisterRoutes(router, authSvc)

	updatesSvc := updates.NewService(updates.NewRepository(pool), logger)
	updates.NewHandler(updatesSvc).RegisterRoutes(router, auth.Middleware(authSvc))

	contextSvc := gcctx.NewService(preferencesSvc, placeSvc, updatesSvc, logger)
	gcctx.NewHandler(contextSvc).RegisterRoutes(router, auth.Middleware(authSvc))

	observationsSvc := observations.NewService(observations.NewRepository(pool), storage.NewMemoryStore(), logger)
	observations.NewHandler(observationsSvc, authSvc).RegisterRoutes(router)

	reportsSvc := reports.NewService(reports.NewRepository(observations.NewRepository(pool)), logger)
	reports.NewHandler(reportsSvc, authSvc).RegisterRoutes(router)

	notifRepo := notifications.NewRepository(pool)
	notifSvc := notifications.NewService(notifRepo, logger)
	notifications.NewHandler(notifSvc).RegisterRoutes(router, auth.Middleware(authSvc))

	auditSvc := audit.NewService(audit.NewRepository(pool))
	audit.NewHandler(auditSvc).RegisterRoutes(router, auth.Middleware(authSvc))

	projectSvc := projects.NewService(projects.NewRepository(pool))
	projects.NewHandler(projectSvc).RegisterRoutes(router, auth.Middleware(authSvc))

	reportingSvc := reporting.NewService(reporting.NewRepository(pool))
	reporting.NewHandler(reportingSvc).RegisterRoutes(router, auth.Middleware(authSvc))

	sourcesSvc := sources.NewService(sources.NewRepository(pool))
	projectListSvc := projects.NewService(projects.NewRepository(pool))
	console.NewHandler(sourcesSvc, projectListSvc).RegisterRoutes(router, auth.Middleware(authSvc))

	integrations.NewHandler("").RegisterRoutes(router)

	server := &http.Server{
		Addr:         net.JoinHostPort(cfg.Server.Host, strconv.Itoa(cfg.Server.Port)),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout),
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout),
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout),
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	logger.Info("http server listening",
		"addr", server.Addr,
		"version", version,
		"app_env", cfg.AppEnv,
	)

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Server.ShutdownTimeout))
		defer cancel()

		logger.Info("shutting down http server")
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("graceful shutdown: %w", err)
		}
		logger.Info("http server stopped cleanly")
		return nil

	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("http server: %w", err)
	}
}
