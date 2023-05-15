package version

import (
	"runtime"
	"runtime/debug"
)

var (
	Version  = "devel"       //nolint:gochecknoglobals
	Revision = "development" //nolint:gochecknoglobals
)

type Report struct {
	Version   string `json:"version,omitempty"`
	Revision  string `json:"revision,omitempty"`
	OS        string `json:"os,omitempty"`
	Arch      string `json:"arch,omitempty"`
	GoVersion string `json:"goVersion,omitempty"`

	Report *debug.BuildInfo `json:"fullReport,omitempty"`
}

func NewReport() *Report {
	goVersion := "unknown"
	revision := ""

	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		goVersion = buildInfo.GoVersion
		// XXX does not work as expected unless go install-ed https://github.com/golang/go/issues/51279
		for _, s := range buildInfo.Settings {
			if s.Key == "vcs.revision" {
				revision = s.Value[:9]
			}
		}
	}

	return &Report{
		Version:   Version,
		Revision:  revision,
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		GoVersion: goVersion,
		Report:    buildInfo,
	}
}
