package filesystem

import "math"

func Init() {
	// Retrieve the current umask so we can restore it on exit
	cMask := umask(0)

	if cMask > math.MaxUint32 || cMask < 0 {
		panic("currently set user umask is out of range")
	}

	currentMask = uint32(cMask)
}
