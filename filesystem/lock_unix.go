// Portions from internal go
//
// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin || dragonfly || freebsd || illumos || linux || netbsd || openbsd

package filesystem

import (
	"errors"
	"os"
	"syscall"

	"go.farcloser.world/core/log"
)

type lockType int16

const (
	readLock  lockType = syscall.LOCK_SH
	writeLock lockType = syscall.LOCK_EX
)

func lock(path string, lockType lockType) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	for {
		err = syscall.Flock(int(file.Fd()), int(lockType))
		if !errors.Is(err, syscall.EINTR) {
			break
		}
	}

	if err != nil {
		if fileErr := file.Close(); fileErr != nil {
			log.Error().Err(fileErr).Msgf("failed closing lock file %q", file.Name())
		}

		return nil, err
	}

	return file, nil
}

func unlock(file *os.File) error {
	defer func() {
		err := file.Close()
		if err != nil {
			log.Error().Err(err).Msgf("failed closing lock file %q", file.Name())
		}
	}()

	for {
		err := syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
		if !errors.Is(err, syscall.EINTR) {
			return err
		}
	}
}
