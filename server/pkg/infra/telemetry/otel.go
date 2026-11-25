// Package telemetry configures and initializes OpenTelemetry providers for observability.
package telemetry

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// InitProvider initializes the OpenTelemetry provider for Traces, Metrics, and Logs.
func InitProvider(ctx context.Context, endpoint, username, apiKey, namespace string, serviceName string) (func(context.Context) error, error) {
	if endpoint == "" {
		return func(context.Context) error { return nil }, nil
	}

	// Set global error handler to use slog at Error level.
	// This ensures OTel internal errors (like connection timeouts) appear as ERROR logs.
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		slog.Error("OpenTelemetry internal error", slog.Any("error", err))
	}))

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			"",
			semconv.ServiceNamespace(namespace),
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+apiKey))
	headers := map[string]string{
		"Authorization": authHeader,
	}

	parsedEndpoint, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid otlp http endpoint: %w", err)
	}
	if parsedEndpoint.Scheme == "" || parsedEndpoint.Host == "" {
		return nil, fmt.Errorf("otlp endpoint must be a full https url, got %q", endpoint)
	}

	traceEndpoint, err := signalEndpoint(endpoint, "traces")
	if err != nil {
		return nil, err
	}
	metricEndpoint, err := signalEndpoint(endpoint, "metrics")
	if err != nil {
		return nil, err
	}
	logEndpoint, err := signalEndpoint(endpoint, "logs")
	if err != nil {
		return nil, err
	}

	// Traces (HTTP)
	traceExporter, err := otlptrace.New(
		ctx,
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpointURL(traceEndpoint),
			otlptracehttp.WithHeaders(headers),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create http trace exporter: %w", err)
	}
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// Metrics (HTTP)
	metricExporter, err := otlpmetrichttp.New(
		ctx,
		otlpmetrichttp.WithEndpointURL(metricEndpoint),
		otlpmetrichttp.WithHeaders(headers),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create http metric exporter: %w", err)
	}
	metricProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter, sdkmetric.WithInterval(30*time.Second))),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(metricProvider)

	// Logs (HTTP)
	logExporter, err := otlploghttp.New(
		ctx,
		otlploghttp.WithEndpointURL(logEndpoint),
		otlploghttp.WithHeaders(headers),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create http log exporter: %w", err)
	}
	logProvider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
		sdklog.WithResource(res),
	)
	global.SetLoggerProvider(logProvider)

	return func(ctx context.Context) error {
		var errs []error
		if err := traceProvider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to shutdown trace provider: %w", err))
		}
		if err := metricProvider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to shutdown metric provider: %w", err))
		}
		if err := logProvider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to shutdown log provider: %w", err))
		}

		if len(errs) > 0 {
			return fmt.Errorf("shutdown errors: %v", errs)
		}
		return nil
	}, nil
}

// signalEndpoint constructs a specific signal endpoint (e.g., traces, metrics) from a generic OTLP base URL.
func signalEndpoint(rawEndpoint, signal string) (string, error) {
	parsed, err := url.Parse(rawEndpoint)
	if err != nil {
		return "", fmt.Errorf("invalid otlp http endpoint: %w", err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("otlp endpoint must be a full https url, got %q", rawEndpoint)
	}

	base := parsed.Scheme + "://" + parsed.Host + strings.TrimSuffix(parsed.Path, "/")
	if strings.HasSuffix(base, "/v1/"+signal) {
		return base, nil
	}
	if strings.Contains(base, "/v1/") {
		before, _, _ := strings.Cut(base, "/v1/")
		return before + "/v1/" + signal, nil
	}
	if strings.HasSuffix(base, "/v1") {
		return base + "/" + signal, nil
	}
	return base + "/v1/" + signal, nil
}
