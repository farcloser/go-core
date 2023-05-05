package attribute

import (
	"go.opentelemetry.io/otel/attribute"
)

func String(k, v string) attribute.KeyValue {
	return attribute.String(k, v)
}
