package telemetry

// traceEndpoint := os.Getenv("OTEL_EXPORTER_JAEGER_ENDPOINT")
// PROMETHEUS ExporterType = "prometheus"
// OTLP       ExporterType = "otlp"

type ExporterType string

type Config struct {
	ServiceName string       `json:"serviceName"`
	Disabled    bool         `json:"disabled"`
	Type        ExporterType `json:"type"`

	// Only for jaegger it seems
	Endpoint string `json:"endpoint"`
}
