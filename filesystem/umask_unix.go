//go:build !windows

package filesystem

import (
	"syscall"
)

func umask(mask int) int {
	return syscall.Umask(mask)
}
