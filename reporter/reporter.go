package reporter

import (
	"net/http"

	"github.com/getsentry/sentry-go"
	"go.farcloser.world/core/log"
	"go.farcloser.world/core/network"
)

// Init should be called when the app starts, from a config object.
func Init(conf *Config) {
	if conf.Disabled {
		log.Warn().Msg("Crash reporting is entirely disabled. This is not recommended.")

		return
	}

	log.Debug().Msg("Initializing crash reporter with config")

	httpClient := &http.Client{}
	if conf.httpClient != nil {
		httpClient = conf.httpClient
	}

	// XXX tricky: this means network MUST be initialized before reporter
	httpClient.Transport = network.GetTransport()

	err := sentry.Init(sentry.ClientOptions{
		HTTPClient:       httpClient,
		Dsn:              conf.DSN,
		Environment:      conf.Environment,
		EnableTracing:    true,
		Release:          conf.Release,
		Debug:            conf.Debug,
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
