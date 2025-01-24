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

package reporter

import (
	"errors"
	"net/http"

	"github.com/getsentry/sentry-go"

	"go.farcloser.world/core/log"
	"go.farcloser.world/core/network"
)

// Init should be called when the app starts, from a config object.
func Init(conf *Config) error {
	if conf.Disabled {
		log.Warn().Msg("Crash reporting is entirely disabled. This is not recommended.")

		return nil
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
		return errors.Join(ErrReporterInitFailed, err)
	}

	return nil
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
