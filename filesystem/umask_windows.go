//go:build windows

package filesystem

func umask(mask int) int {
	return 0
}
