//go:build windows

package config

func umask(mask int) int {
	return 0
}
