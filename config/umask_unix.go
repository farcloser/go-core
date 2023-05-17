//go:build !windows

package config

import (
	"syscall"
)

func umask(mask int) int {
	return syscall.Umask(mask)
}
