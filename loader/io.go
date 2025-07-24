/*
   Copyright Farcloser.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

//revive:disable:confusing-naming

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

//nolint:wrapcheck
func read(cfg any, location ...string) error {
	loc := absolute(location...)

	mut.Lock()
	defer mut.Unlock()

	//nolint:gosec
	data, err := os.ReadFile(loc)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &cfg)
}

func write(cfg any, location ...string) error {
	loc := absolute(location...)

	mut.Lock()
	defer mut.Unlock()

	err := os.MkdirAll(path.Dir(loc), filesystem.DirPermissionsDefault)
	if err != nil {
		//nolint:wrapcheck
		return err
	}

	data, err := json.MarshalIndent(&cfg, "", " ")
	if err != nil {
		//nolint:wrapcheck
		return err
	}

	//nolint:wrapcheck
	return filesystem.WriteFile(loc, data, filesystem.FilePermissionsDefault)
}

func remove(location ...string) error {
	loc := absolute(location...)

	mut.Lock()
	defer mut.Unlock()

	//nolint:wrapcheck
	return os.Remove(loc)
}
