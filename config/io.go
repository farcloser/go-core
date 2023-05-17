package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"

	"go.codecomet.dev/core/filesystem"
)

var mut *sync.Mutex //nolint:gochecknoglobals

func Absolute(location ...string) string {
	loc := path.Join(location...)

	if !filepath.IsAbs(loc) {
		dir, _ := os.UserConfigDir()
		loc = path.Join(dir, loc)
	}

	return loc
}

func Read(cfg interface{}, location ...string) error {
	loc := Absolute(location...)

	if mut == nil {
		mut = &sync.Mutex{}
	}

	mut.Lock()
	defer mut.Unlock()

	data, err := os.ReadFile(loc)
	if err != nil {
		return fmt.Errorf("failed reading config file %w", err)
	}

	return json.Unmarshal(data, &cfg)
}

func Write(cfg interface{}, location ...string) error {
	loc := Absolute(location...)

	if mut == nil {
		mut = &sync.Mutex{}
	}

	mut.Lock()
	defer mut.Unlock()

	err := os.MkdirAll(path.Dir(loc), defaultDirPerms)
	if err != nil {
		return fmt.Errorf("failed creating config parent directory %w", err)
	}

	data, err := json.MarshalIndent(&cfg, "", " ")
	if err != nil {
		return fmt.Errorf("failed marshalling config json %w", err)
	}

	return filesystem.WriteFile(loc, data, defaultFilePerms)
}

// Delete destroys the config file.
func Delete(location ...string) error {
	loc := Absolute(location...)

	if mut == nil {
		mut = &sync.Mutex{}
	}

	mut.Lock()
	defer mut.Unlock()

	return os.Remove(loc)
}
