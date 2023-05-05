package exec

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"go.codecomet.dev/core/log"
	"go.codecomet.dev/core/reporter"
)

func Resolve(bin string) (string, error) {
	o, err := exec.Command("which", bin).Output()
	out := string(o)
	out = strings.Trim(out, "\n")

	if err != nil {
		err = fmt.Errorf("Resolve errored with: %w", err)
	}

	return out, err
}

func New(defaultBin string, envBin string) *Commander {
	// This is only useful for test...
	bin := os.Getenv(envBin)
	if bin == "" {
		bin = defaultBin
	}

	execut := bin
	// XXX this is ill-designed
	if !filepath.IsAbs(bin) {
		var err error
		execut, err = os.Executable()

		if err != nil {
			reporter.CaptureException(fmt.Errorf("failed retrieving current binary information: %w", err))
			log.Fatal().Err(err).Msg("Cannot find current binary location. This is very wrong.")
		}

		execut = filepath.Join(filepath.Dir(execut), bin)

		if _, err := os.Stat(execut); err != nil {
			// Fallback to path resolution
			execut, _ = Resolve(bin)
		}
	}

	if _, err := os.Stat(execut); err != nil {
		w, _ := os.Getwd()
		reporter.CaptureException(fmt.Errorf("failed finding cli %s with pwd %s - err: %w", bin, w, err))
		log.Fatal().Str("pwd", w).Msgf("Failed finding cli %s with pwd %s - err: %s", bin, w, err)
	}

	return &Commander{
		mu:  &sync.Mutex{},
		bin: execut,
	}
}

type Commander struct {
	mu       *sync.Mutex
	bin      string
	Stdin    io.Reader
	Env      map[string]string
	PreArgs  []string
	Dir      string
	NoReport bool
}

func (com *Commander) Attach(args ...string) error {
	var err error
	if com.Stdin != nil {
		err = com.ExecPipes(com.Stdin, os.Stdout, os.Stderr, args...)
	} else {
		err = com.ExecPipes(os.Stdin, os.Stdout, os.Stderr, args...)
	}

	if err != nil && !com.NoReport {
		reporter.CaptureException(fmt.Errorf("failed attached execution: %w", err))
		log.Error().Err(err).Msg("Attached execution failed")
	}

	return err
}

func (com *Commander) Exec(args ...string) (string, string, error) {
	var stdout bytes.Buffer

	var stderr bytes.Buffer

	err := com.ExecPipes(com.Stdin, &stdout, &stderr, args...)
	sout := stdout.String()
	serr := stderr.String()

	if !com.NoReport && err != nil {
		reporter.CaptureException(fmt.Errorf("failed sub execution: %w - out: %s - err: %s", err, sout, serr))
		log.Error().Err(err).Str("out", sout).Str("err", serr).Msg("Execution failed")
	}

	return sout, serr, err
}

func (com *Commander) ExecPipes(stdin io.Reader, stdout io.Writer, stderr io.Writer, args ...string) error {
	args = append(com.PreArgs, args...)

	envs := []string{}
	for k, v := range com.Env {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}

	log.Debug().Str("binary", com.bin).Strs("arguments", args).Strs("env", envs).Msg("Executing command")

	command := exec.Command(com.bin, args...) //nolint:gosec
	if com.Dir != "" {
		command.Dir = com.Dir
	}

	command.Env = append(os.Environ(), envs...)

	command.Stdin = stdin
	command.Stdout = stdout
	command.Stderr = stderr

	com.mu.Lock()
	err := command.Run()
	com.mu.Unlock()

	if err != nil {
		err = fmt.Errorf("ExecPipes errored: %w", err)
	}

	return err
}
