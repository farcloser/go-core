//go:build !windows

package config

import (
	"syscall"
)

func umask(mask int) {
	syscall.Umask(mask)
}
