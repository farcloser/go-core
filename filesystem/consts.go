package filesystem

const (
	FilePermissionsDefault = 0o644
	DirPermissionsDefault  = 0o755
	FilePermissionsPrivate = 0o600
	DirPermissionsPrivate  = 0o700

	defaultUmask = 0o077
)
