// From go internal
//
// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

// Package filelock provides a platform-independent API for advisory file
// locking. Calls to functions in this package on platforms that do not support
// advisory locks will return errors for which IsNotSupported returns true.

import (
	"errors"
	"os"
)

// Lock places an advisory write lock on the file, blocking until it can be
// locked.
//
// If Lock returns nil, no other process will be able to place a read or write
// lock on the file until this process exits, closes f, or calls Unlock on it.
func Lock(path string) (*os.File, error) {
	file, err := lock(path, writeLock)
	if err != nil {
		err = errors.Join(ErrLockFail, err)
	}

	return file, err
}

// ReadOnlyLock places an advisory read lock on the file, blocking until it can be locked.
//
// If ReadOnlyLock returns nil, no other process will be able to place a write lock on
// the file until this process exits, closes f, or calls Unlock on it.
func ReadOnlyLock(path string) (*os.File, error) {
	file, err := lock(path, readLock)
	if err != nil {
		err = errors.Join(ErrLockFail, err)
	}

	return file, err
}

// Unlock removes an advisory lock placed on f by this process.
func Unlock(lock *os.File) error {
	if lock == nil {
		return ErrLockIsNil
	}

	err := unlock(lock)
	if err != nil {
		err = errors.Join(ErrUnlockFail, err)
	}

	return err
}

func WithLock(path string, function func() error) (err error) {
	file, err := Lock(path)
	if err != nil {
		return err
	}

	defer func() {
		err = errors.Join(Unlock(file), err)
	}()

	return function()
}

func WithReadOnlyLock(path string, function func() error) (err error) {
	file, err := ReadOnlyLock(path)
	if err != nil {
		return err
	}

	defer func() {
		err = errors.Join(Unlock(file), err)
	}()

	return function()
}
