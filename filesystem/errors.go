package filesystem

import "errors"

var (
	ErrLockFail        = errors.New("failed to acquire lock")
	ErrUnlockFail      = errors.New("failed to release lock")
	ErrAtomicWriteFail = errors.New("failed to write file atomically")
	ErrLockIsNil       = errors.New("nil lock")
	ErrInvalidPath     = errors.New("invalid path")

	errInvalidPathTooLong = errors.New("path component must be stricly shorter than 256 characters")
	errInvalidPathEmpty   = errors.New("path component cannot be empty")
)
