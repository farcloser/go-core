// From internal go
//
// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// https://cs.opensource.google/go/go/+/master:src/cmd/go/internal/lockedfile/internal/filelock/filelock_windows.go

//go:build windows

package filesystem

import (
	"errors"
	"os"

	"golang.org/x/sys/windows"
)

type lockType uint32

const (
	// https://msdn.microsoft.com/en-us/library/windows/desktop/aa365203(v=vs.85).aspx
	readLock  lockType = 0
	writeLock lockType = windows.LOCKFILE_EXCLUSIVE_LOCK

	reserved = 0
	allBytes = ^uint32(0)

	lockPermission = 0o600
)

func lock(path string, lockType lockType) (file *os.File, err error) {
	file, err = os.OpenFile(path+".lock", os.O_CREATE, lockPermission)
	if err != nil {
		return nil, err
	}

	if err = windows.LockFileEx(
		windows.Handle(file.Fd()),
		uint32(lockType), reserved, allBytes, allBytes, new(windows.Overlapped)); err != nil {
		if fileErr := file.Close(); fileErr != nil {
			err = errors.Join(err, fileErr)
		}

		return nil, err
	}

	return file, nil
}

func unlock(file *os.File) (err error) {
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	return windows.UnlockFileEx(windows.Handle(file.Fd()), reserved, allBytes, allBytes, new(windows.Overlapped))
}
