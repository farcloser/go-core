//go:build windows

package filesystem

func umask(_ int) int {
	return 0
}
