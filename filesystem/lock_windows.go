// From internal go
//
// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package filesystem

import (
	"errors"
	"os"
	"syscall"
)

type lockType uint32

const (
	// https://msdn.microsoft.com/en-us/library/windows/desktop/aa365203(v=vs.85).aspx
	readLock  lockType = 0
	writeLock lockType = syscall.LOCKFILE_EXCLUSIVE_LOCK
)

const (
	reserved = 0
	allBytes = ^uint32(0)
)

func lock(path string, lockType lockType) (file *os.File, err error) {
	file, err = os.OpenFile(path+".lock", os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	if err = syscall.LockFileEx(syscall.Handle(f.Fd()), uint32(lockType), reserved, allBytes, allBytes, new(syscall.Overlapped)); err != nil {
		if fileErr := file.Close(); fileErr != nil {
			err = errors.Join(err, fileErr)
		}

		return nil, err
	}

	return file, nil
}

func unlock(locked File) (err error) {
	defer func() {
		if closeErr = locked.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	return syscall.UnlockFileEx(syscall.Handle(f.Fd()), reserved, allBytes, allBytes, new(syscall.Overlapped))
}
