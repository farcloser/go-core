//nolint:ireturn
package telemetry

import (
	"context"
	"fmt"
	"io"
	"time"

	sentryotel "github.com/getsentry/sentry-go/otel"
	"go.farcloser.world/core/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

const closeTimeout = 5 * time.Second

type TracerProvider = trace.TracerProvider

func GetTracerProvider() TracerProvider {
	return otel.GetTracerProvider()
}

func Init(conf *Config) io.Closer {
	if conf.Disabled {
		log.Warn().Msg("Telemetry is disabled.")

		return &noopCloser{}
	}

	prov, err := provider(conf.Type, conf.Endpoint, conf.ServiceName)
	if err != nil {
		log.Fatal().Err(err).Str("type", string(conf.Type)).Msg("Failed creating telemetry provider")
	}

	// Register with OTEL
	otel.SetTracerProvider(prov)

	return providerCloser{
		TracerProvider: prov,
	}
}

type noopCloser struct{}

func (*noopCloser) Close() error {
	return nil
}

type providerCloser struct {
	*sdktrace.TracerProvider
}

func (t providerCloser) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), closeTimeout)
	defer cancel()

	return t.Shutdown(ctx)
}

func provider(expType ExporterType, url string, serviceName string) (*sdktrace.TracerProvider, error) {
	var err error

	var exp sdktrace.SpanExporter

	opts := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	}

	switch expType {
	case JAEGGER:
		exp, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
		opts = append(opts, sdktrace.WithBatcher(exp, sdktrace.WithMaxExportBatchSize(1)))
	case SENTRY:
		opts = append(opts, sdktrace.WithSpanProcessor(sentryotel.NewSentrySpanProcessor()))
		otel.SetTextMapPropagator(sentryotel.NewSentryPropagator())
	/*
		case PROMETHEUS:
		case OTLP:

	*/
	default:
		err = ErrUnsupportedProviderType
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		opts...,
	)

	return tracerProvider, nil
}
