package loader

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"sync"

	"go.farcloser.world/core/filesystem"
)

var mut sync.Mutex //nolint:gochecknoglobals

func absolute(location ...string) string {
	loc := path.Join(location...)

	if !filepath.IsAbs(loc) {
		dir, _ := os.UserConfigDir()
		loc = path.Join(dir, loc)
	}

	return loc
}

func read(cfg interface{}, location ...string) error {
	loc := absolute(location...)

	mut.Lock()
	defer mut.Unlock()

	data, err := os.ReadFile(loc)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &cfg)
}

func write(cfg interface{}, location ...string) error {
	loc := absolute(location...)

	mut.Lock()
	defer mut.Unlock()

	err := os.MkdirAll(path.Dir(loc), filesystem.DirPermissionsDefault)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(&cfg, "", " ")
	if err != nil {
		return err
	}

	return filesystem.WriteFile(loc, data, filesystem.FilePermissionsDefault)
}

func remove(location ...string) error {
	loc := absolute(location...)

	mut.Lock()
	defer mut.Unlock()

	return os.Remove(loc)
}
