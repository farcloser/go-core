package filesystem

var currentMask = defaultUmask //nolint:gochecknoglobals

func SetUmask(mask int) {
	if mask == currentMask {
		return
	}

	currentMask = mask
	umask(mask)
}

func GetUmask() int {
	return currentMask
}
