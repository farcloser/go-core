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

package log

import (
	"bufio"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init should be called when the app starts, from a config object.
func Init(conf *Config) {
	// This mostly should be the responsibility of the app itself but hey
	zerolog.SetGlobalLevel(conf.Level)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	zerolog.TimestampFunc = func() time.Time {
		// XXX investigate this
		//nolint:gosmopolitan
		return time.Now().In(time.Local)
	}
}

// SetLevel sets the global log level.
func SetLevel(lv Level) {
	zerolog.SetGlobalLevel(lv)
}

// GetLevel returns the current global log level.
func GetLevel() Level {
	return zerolog.GlobalLevel()
}

// DebugSink is a sink for debug logs that reads from the provided reader.
func DebugSink(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		log.Debug().Msg(scanner.Text())
	}
}

// WarnSink is a sink for warning logs that reads from the provided reader.
func WarnSink(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		log.Warn().Msg(scanner.Text())
	}
}

// ErrorSink is a sink for error logs that reads from the provided reader.
func ErrorSink(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		log.Error().Msg(scanner.Text())
	}
}

// Error returns an Event for an error.
//
//nolint:zerologlint
func Error() *Event {
	return log.Error()
}

// Warn returns an Event for a warning.
//
//nolint:zerologlint
func Warn() *Event {
	return log.Warn()
}

// Info returns an Event for informational messages.
//
//nolint:zerologlint
func Info() *Event {
	return log.Info()
}

// Debug returns an Event for debug messages.
//
//nolint:zerologlint
func Debug() *Event {
	return log.Debug()
}

// Fatal returns an Event for fatal messages and exits the application.
//
//nolint:zerologlint
func Fatal() *Event {
	return log.Fatal()
}
