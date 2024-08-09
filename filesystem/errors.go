package filesystem

import "errors"

var (
	ErrLockFail        = errors.New("failed to acquire lock")
	ErrUnlockFail      = errors.New("failed to release lock")
	ErrAtomicWriteFail = errors.New("failed to write file atomically")
)
