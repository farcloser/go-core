package reporter

import (
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"go.codecomet.dev/core/log"
	"go.codecomet.dev/core/network"
)

const flushTimeout = 2 * time.Second

// Init should be called when the app starts, from a config object.
func Init(cnf *Config) {
	if cnf.Disabled {
		log.Warn().Msg("Crash reporting is entirely disabled. This is not recommended.")

		return
	}

	log.Debug().Msg("Initializing crash reporter with config")

	httpClient := cnf.httpClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	httpClient.Transport = network.Get().Transport()

	err := sentry.Init(sentry.ClientOptions{
		HTTPClient:       httpClient,
		Dsn:              cnf.DSN,
		Environment:      cnf.Environment,
		EnableTracing:    true,
		Release:          cnf.Release,
		Debug:            cnf.Debug,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("sentry.Init failed")
	}
}

func CaptureException(err error) *EventID {
	return sentry.CaptureException(err)
}

func CaptureMessage(msg string) *EventID {
	return sentry.CaptureMessage(msg)
}

func CaptureEvent(e *Event) *EventID {
	return sentry.CaptureEvent(e)
}

func Shutdown() {
	// Flush buffered events before the program terminates.
	// Set the timeout to the maximum duration the program can afford to wait.
	sentry.Flush(flushTimeout)
}
