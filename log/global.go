package log

import (
	"bufio"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"time"
)

// Init should be called when the app starts, from a config object
func Init(cnf *Config) {
	// This mostly should be the responsibility of the app itself but hey
	zerolog.SetGlobalLevel(cnf.Level)
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

func DebugSink(o string) {
	scanner := bufio.NewScanner(strings.NewReader(o))
	for scanner.Scan() {
		log.Debug().Msg(scanner.Text())
	}
}

func WarnSink(o string) {
	scanner := bufio.NewScanner(strings.NewReader(o))
	for scanner.Scan() {
		log.Warn().Msg(scanner.Text())
	}
}

func ErrorSink(o string) {
	scanner := bufio.NewScanner(strings.NewReader(o))
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
