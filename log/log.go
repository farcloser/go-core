//nolint:zerologlint
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
		return time.Now().In(time.Local)
	}
}

func SetLevel(lv Level) {
	zerolog.SetGlobalLevel(lv)
}

func GetLevel() Level {
	return zerolog.GlobalLevel()
}

func DebugSink(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		log.Debug().Msg(scanner.Text())
	}
}

func WarnSink(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		log.Warn().Msg(scanner.Text())
	}
}

func ErrorSink(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		log.Error().Msg(scanner.Text())
	}
}

func Error() *Event {
	return log.Error()
}

func Warn() *Event {
	return log.Warn()
}

func Info() *Event {
	return log.Info()
}

func Debug() *Event {
	return log.Debug()
}

func Fatal() *Event {
	return log.Fatal()
}
