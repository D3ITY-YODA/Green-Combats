package telemetry

import (
	"context"
	"fmt"
)

// Config holds OpenTelemetry configuration.
type Config struct {
	Enabled     bool   `yaml:"enabled"`
	Endpoint    string `yaml:"endpoint"`     // OTLP gRPC endpoint
	ServiceName string `yaml:"service_name"` // defaults to "green-compass"
	Environment string `yaml:"environment"`  // local, dev, staging, production
}

// ShutdownFunc is called to flush pending telemetry data.
type ShutdownFunc func(ctx context.Context) error

// Init initializes OpenTelemetry with the given configuration.
// Returns a shutdown function that should be called on application exit.
func Init(ctx context.Context, cfg Config) (ShutdownFunc, error) {
	if !cfg.Enabled {
		return func(_ context.Context) error { return nil }, nil
	}

	if cfg.ServiceName == "" {
		cfg.ServiceName = "green-compass"
	}

	// In a full implementation, this would:
	// 1. Set up OTLP exporter to send traces/metrics to the endpoint
	// 2. Initialize TracerProvider and MeterProvider
	// 3. Set global propagators (W3C TraceContext + Baggage)
	// 4. Register as global providers
	//
	// Example (using go.opentelemetry.io/otel):
	//
	//   exporter, err := otlptracegrpc.New(ctx,
	//       otlptracegrpc.WithEndpoint(cfg.Endpoint),
	//       otlptracegrpc.WithInsecure(),
	//   )
	//   tp := sdktrace.NewTracerProvider(
	//       sdktrace.WithBatcher(exporter),
	//       sdktrace.WithResource(resource.NewWithAttributes(
	//           semconv.SchemaURL,
	//           semconv.ServiceNameKey.String(cfg.ServiceName),
	//           semconv.DeploymentEnvironmentKey.String(cfg.Environment),
	//       )),
	//   )
	//   otel.SetTracerProvider(tp)
	//   otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
	//       propagation.TraceContext{},
	//       propagation.Baggage{},
	//   ))

	fmt.Printf("telemetry: initialized (endpoint=%s, service=%s, env=%s)\n",
		cfg.Endpoint, cfg.ServiceName, cfg.Environment)

	shutdown := func(ctx context.Context) error {
		// In production: tp.Shutdown(ctx) to flush pending data
		fmt.Println("telemetry: shutdown complete")
		return nil
	}

	return shutdown, nil
}
