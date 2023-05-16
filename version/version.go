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
