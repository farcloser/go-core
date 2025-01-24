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

package exec

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"go.farcloser.world/core/log"
	"go.farcloser.world/core/reporter"
)

// FIXME: rewrite this whole stuff

var ErrExecResolutionFail = errors.New("resolve errored")

func Resolve(binaryName string) (string, error) {
	// https://unix.stackexchange.com/questions/85249/why-not-use-which-what-to-use-then
	o, err := exec.Command("command", "-v", binaryName).Output()
	if err != nil {
		return "", errors.Join(ErrExecResolutionFail, err)
	}

	return strings.Trim(string(o), "\n"), nil
}

func New(defaultBinaryLocation string, environVariable string) *Commander {
	// This is only useful for integration testing, really...
	binaryLocation := os.Getenv(environVariable)
	if binaryLocation == "" {
		binaryLocation = defaultBinaryLocation
	}

	resolvedLocation := binaryLocation
	// XXX this is ill-designed
	if !filepath.IsAbs(binaryLocation) {
		var err error

		resolvedLocation, err = os.Executable()
		if err != nil {
			reporter.CaptureException(fmt.Errorf("failed retrieving current binary information: %w", err))
			log.Fatal().Err(err).Msg("Cannot find current binary location. This is very wrong.")
		}

		resolvedLocation = filepath.Join(filepath.Dir(resolvedLocation), binaryLocation)

		if _, err = os.Stat(resolvedLocation); err != nil {
			// Fallback to path resolution
			resolvedLocation, _ = Resolve(binaryLocation)
		}
	}

	if _, err := os.Stat(resolvedLocation); err != nil {
		w, _ := os.Getwd()
		reporter.CaptureException(fmt.Errorf("failed finding cli %s with pwd %s - err: %w", binaryLocation, w, err))
		log.Fatal().Str("pwd", w).Msgf("Failed finding cli %s with pwd %s - err: %s", binaryLocation, w, err)
	}

	return &Commander{
		mu:  &sync.Mutex{},
		bin: resolvedLocation,
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
