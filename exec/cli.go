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

	"github.com/codecomet-io/go-core/log"
	"github.com/codecomet-io/go-core/reporter"
)

func Resolve(bin string) (string, error) {
	o, err := exec.Command("which", bin).Output()
	out := string(o)
	out = strings.Trim(out, "\n")

	return out, err
}

func New(defaultBin string, envBin string) *Commander {
	// This is only useful for test...
	bin := os.Getenv(envBin)
	if bin == "" {
		bin = defaultBin
	}

	ex := bin
	// XXX this is ill-designed
	if !filepath.IsAbs(bin) {
		var err error
		ex, err = os.Executable()
		if err != nil {
			reporter.CaptureException(fmt.Errorf("failed retrieving current binary information: %s", err))
			log.Fatal().Err(err).Msg("Cannot find current binary location. This is very wrong.")
		}
		ex = filepath.Join(filepath.Dir(ex), bin)

		if _, err := os.Stat(ex); err != nil {
			// Fallback to path resolution
			ex, _ = Resolve(bin)
		}
	}

	if _, err := os.Stat(ex); err != nil {
		w, _ := os.Getwd()
		reporter.CaptureException(fmt.Errorf("failed finding cli %s with pwd %s - err: %s", bin, w, err))
		log.Fatal().Str("pwd", w).Msgf("Failed finding cli %s with pwd %s - err: %s", bin, w, err)
	}

	return &Commander{
		mu:  &sync.Mutex{},
		bin: ex,
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
		reporter.CaptureException(fmt.Errorf("failed attached execution: %s", err))
		log.Error().Err(err).Msg("Attached execution failed")
	}

	return err
}

func (com *Commander) Exec(args ...string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := com.ExecPipes(com.Stdin, &stdout, &stderr, args...)
	o := stdout.String()
	e := stderr.String()
	if err != nil && !com.NoReport {
		reporter.CaptureException(fmt.Errorf("failed sub execution: %s - out: %s - err: %s", err, o, e))
		log.Error().Err(err).Str("out", o).Str("err", e).Msg("Execution failed")
	}

	return o, e, err
}

func (com *Commander) ExecPipes(stdin io.Reader, stdout io.Writer, stderr io.Writer, args ...string) error {
	args = append(com.PreArgs, args...)

	var envs []string
	for k, v := range com.Env {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}
	log.Debug().Str("binary", com.bin).Strs("arguments", args).Strs("env", envs).Msg("Executing command")

	c := exec.Command(com.bin, args...)
	if com.Dir != "" {
		c.Dir = com.Dir
	}
	c.Env = append(os.Environ(), envs...)
	c.Stdin = stdin
	c.Stdout = stdout
	c.Stderr = stderr

	com.mu.Lock()
	e := c.Run()
	com.mu.Unlock()

	return e
}
