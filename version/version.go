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

package version

import (
	"runtime"
	"runtime/debug"
)

const unknown = "unknown"

var Version = unknown //nolint:gochecknoglobals

type Report struct {
	Version   string `json:"version,omitempty"`
	Revision  string `json:"revision,omitempty"`
	Dirty     bool   `json:"dirty,omitempty"`
	OS        string `json:"os,omitempty"`
	Arch      string `json:"arch,omitempty"`
	GoVersion string `json:"goVersion,omitempty"`

	Raw *debug.BuildInfo `json:"rawReport,omitempty"`
}

func NewReport() *Report {
	rep := &Report{
		Version:   Version,
		Revision:  unknown,
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		GoVersion: unknown,
		Dirty:     false,
	}

	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		rep.Raw = buildInfo
		rep.GoVersion = buildInfo.GoVersion
		// XXX is this really working as expected? may depend on go version...
		// unless go install-ed https://github.com/golang/go/issues/51279
		for _, s := range buildInfo.Settings {
			if s.Key == "vcs.revision" {
				rep.Revision = s.Value[:7]
			}

			if s.Key == "vcs.modified" && s.Value == "true" {
				rep.Dirty = true
			}
		}
	}

	return rep
}
