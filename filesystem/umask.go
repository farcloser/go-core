package filesystem

var currentMask = defaultUmask //nolint:gochecknoglobals

func SetUmask(mask uint32) {
	if mask == currentMask {
		return
	}

	currentMask = mask
	_ = umask(int(mask))
}

func GetUmask() uint32 {
	return currentMask
}
