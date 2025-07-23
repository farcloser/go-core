/*
   Copyright Farcloser.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package telemetry

import (
	"context"
	"errors"
	"io"
	"time"

	sentryotel "github.com/getsentry/sentry-go/otel"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"

	"go.farcloser.world/core/log"
	"go.farcloser.world/core/network"
)

const closeTimeout = 5 * time.Second

// TracerProvider provides Tracers that are used by instrumentation code to
// trace computational workflows.
type TracerProvider = trace.TracerProvider

type noopCloser struct{}

func (*noopCloser) Close() error {
	return nil
}

type providerCloser struct {
	*sdktrace.TracerProvider
}

// Close closes the tracer provider and any associated resources.
func (t providerCloser) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), closeTimeout)
	defer cancel()

	err := t.Shutdown(ctx)
	if err != nil {
		err = errors.Join(ErrCloseError, err)
	}

	return err
}

// GetTracerProvider returns the registered global trace provider.
//
//nolint:ireturn
func GetTracerProvider() TracerProvider {
	return otel.GetTracerProvider()
}

// Init initializes the telemetry provider based on the provided configuration.
func Init(conf *Config) (io.Closer, error) {
	if conf.Disabled {
		log.Warn().Msg("Telemetry is disabled.")

		return &noopCloser{}, nil
	}

	prov, err := provider(conf.Type, conf.Endpoint, conf.ServiceName)
	if err != nil {
		return nil, err
	}

	// Register with OTEL
	otel.SetTracerProvider(prov)

	return providerCloser{
		TracerProvider: prov,
	}, nil
}

func provider(expType ExporterType, url, serviceName string) (*sdktrace.TracerProvider, error) {
	var err error

	var exp sdktrace.SpanExporter

	opts := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	}

	switch expType {
	case JAEGER:
		ctx := context.Background()

		exp, err = otlptracehttp.New(
			ctx,
			otlptracehttp.WithEndpoint(url),
			otlptracehttp.WithTLSClientConfig(network.GetTLSConfig()),
		)
		if err != nil {
			panic(err)
		}

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
		return nil, errors.Join(ErrProviderCreationFailed, err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		opts...,
	)

	return tracerProvider, nil
}
