package term

import "github.com/mattn/go-isatty"

func IsTerminal(fd uintptr) bool {
	return isatty.IsTerminal(fd)
}
