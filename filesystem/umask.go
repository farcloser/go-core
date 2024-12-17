package filesystem

var currentMask = defaultUmask //nolint:gochecknoglobals

func SetUmask(mask uint32) {
	if mask == currentMask {
		return
	}

	currentMask = mask
	umask(int(mask)) //nolint:staticcheck
}

func GetUmask() uint32 {
	return currentMask
}
