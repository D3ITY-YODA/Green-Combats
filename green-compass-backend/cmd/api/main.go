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

	"green-compass-backend/internal/config"
	"green-compass-backend/internal/health"
	"green-compass-backend/pkg/clock"
	"green-compass-backend/pkg/httpx"
	"green-compass-backend/pkg/logging"
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

	router := gin.New()
	router.Use(gin.Recovery(), httpx.RequestLogger(logger, clock.New()))

	healthService := health.NewService(version, clock.New())
	health.NewHandler(healthService).RegisterRoutes(router)

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
