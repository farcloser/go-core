package filesystem

func Init() {
	currentMask = umask(0)
	umask(currentMask)
}
