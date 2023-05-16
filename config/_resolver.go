package config

/*

var values map[string]string

type resolver struct {
	defaultValue string
	envKey       string
}

var resolvers map[string]*resolver

func Resolve(key string) (string, bool) {
	// If already resolved, return it
	if val, ok := values[key]; ok {
		return val, true
	}

	// See if we have a resolver, and if so, resolve, store and return
	if resolver, ok := resolvers[key]; ok {
		values[key] = resolver.defaultValue
		if resolver.envKey != "" {
			if ret, ok := os.LookupEnv(resolver.envKey); ok {
				values[key] = ret
			}
		}
		return values[key], true
	}
	return "", false
}

func Declare(key string, envKey string, defaultValue string) {
	resolvers[key] = &resolver{
		defaultValue: defaultValue,
		envKey:       envKey,
	}
}

Base.Root = Resolve("CODECOMET_ROOT", "/var/lib/isovaline")

// Base describes a generic configuration object, with path creation / abs. logic and read write facility.
type Base struct {
	appName string
	root    string

	// DataRoot  string      `json:"dataRoot,omitempty"`
	Location  []string    `json:"-"`

	mu *sync.Mutex
}

// GetDataRoot returns the root directory for persistent data storage.
func (base *Base) GetDataRoot() string {
	if defaults.DataRootEnv != "" {
		return defaults.DataRootEnv
	}
	return defaults.DataRoot
}

// GetCacheDir returns the root directory for temporary files.
func (base *Base) GetCacheDir() string {
	dir, _ := os.UserCacheDir()
	return path.Join(dir, base.appName)
}

func (base *Base) GetConfigDir() string {
	dir, _ := os.UserConfigDir()
	return path.Join(dir, base.appName)
}

func (base *Base) GetHomeDir() string {
	dir, _ := os.UserHomeDir()
	return dir
}

func (base *Base) GetTempDir() string {
	return os.TempDir()
}

// ensureExistence guarantees that a directory exists, possibly creating it with dir permissions (unsafe?)
func (base *Base) ensureExistence(pth string) error {
	return os.MkdirAll(pth, base.DirPerms)
}

// getLocation returns the config file location, accounting for env and absolutization.
func (base *Base) getLocation() string {
	loc := getAbsolute(path.Join(base.Location...), configRoot())
	err := base.ensureExistence(filepath.Dir(loc))
	if err != nil {
		panic(err)
	}
	return loc
}


*/
