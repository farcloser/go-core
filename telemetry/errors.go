package telemetry

import "errors"

var (
	ErrUnsupportedProviderType = errors.New("unsupported provider type")
	ErrProviderCreationFailed  = errors.New("provider creation failed")
)
