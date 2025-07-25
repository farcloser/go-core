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

package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"go.farcloser.world/core/filesystem"
	"go.farcloser.world/core/log"
	"go.farcloser.world/core/network"
	"go.farcloser.world/core/reporter"
	"go.farcloser.world/core/telemetry"
)

// New creates a new Core configuration object.
func New(appName string, location ...string) *Core {
	conf := &Core{
		location: append([]string{appName}, location...),

		Client: &network.Config{
			TLSMin:              defaultTLSClientMinVersion,
			TLSHandshakeTimeout: defaultTLSHandshakeTimeout,
			DialerKeepAlive:     defaultDialerKeepAlive,
			DialerTimeout:       defaultDialerTimeout,
			DisallowSystemRoot:  false,
			CertPath:            defaultCertPath,
			KeyPath:             defaultKeyPath,
			RootCAs:             []string{},
		},

		Server: &network.Config{
			TLSMin:              defaultTLSServerMinVersion,
			TLSHandshakeTimeout: defaultTLSHandshakeTimeout,
			CertPath:            defaultCertPath,
			KeyPath:             defaultKeyPath,
		},

		Logger: &log.Config{
			Level: defaultLogLevel,
		},
	}

	conf.Client.Resolve = conf.Resolve
	conf.Server.Resolve = conf.Resolve

	return conf
}

// Core is the core configuration object.
type Core struct {
	Reporter  *reporter.Config  `json:"reporter,omitempty"`
	Logger    *log.Config       `json:"logger,omitempty"`
	Telemetry *telemetry.Config `json:"telemetry,omitempty"`
	Client    *network.Config   `json:"client,omitempty"`
	Server    *network.Config   `json:"server,omitempty"`

	Umask uint32 `json:"umask,omitempty"`

	location []string
}

// Trust does trust a certificate for both client and server.
func (obj *Core) Trust(ca ...string) {
	if len(ca) > 0 {
		obj.Server.ClientCA = ca[0]
		obj.Client.RootCAs = ca
	}
}

// Resolve resolves a location against the config file's directory.
func (obj *Core) Resolve(location ...string) string {
	// Get the absolute path of the containing dir of the config file, resolved against UserConfigDir
	base := absolute(obj.location[:len(obj.location)-1]...)

	// If the desired location is not absolute, resolve against above
	loc := path.Join(location...)
	if !filepath.IsAbs(loc) {
		loc = path.Join(append([]string{base}, location...)...)
	}

	return loc
}

// Ensure ensures that the parent directory of the given location exists.
func (obj *Core) Ensure(location ...string) error {
	// Get the absolute path of the containing dir of the config file, resolved against UserConfigDir
	base := absolute(obj.location[:len(obj.location)-1]...)

	loc := path.Join(location...)
	if !filepath.IsAbs(loc) {
		loc = path.Join(append([]string{base}, location...)...)
	}

	err := os.MkdirAll(path.Dir(loc), filesystem.DirPermissionsDefault)
	if err != nil {
		err = fmt.Errorf("failed to ensure parent directory existence for %s: %w", loc, err)
	}

	return err
}

// OnIO is called when the IO subsystem is initialized.
func (obj *Core) OnIO() {
	// Note: calling init everytime we load is not super efficient, but then, how often does that happen?
	// Init filesystem first (capture the current, actual umask before we do anything)
	filesystem.Init()
	// Now, set the umask to whatever
	filesystem.SetUmask(obj.Umask)
}

// GetLocation returns the location of the config file.
func (obj *Core) GetLocation() []string {
	return obj.location
}

// GetDataRoot returns the data root directory for the application.
func (obj *Core) GetDataRoot() string {
	var loc string

	base, _ := os.UserHomeDir()

	switch runtime.GOOS {
	case "darwin":
		// XXX figure out impact on iCloud auto backup thing and containers
		loc = path.Join(base, "Library", "Application Support", obj.location[0])
	default:
		loc = path.Join(base, "."+obj.location[0])
	}

	// XXX ignore errors?
	_ = os.MkdirAll(path.Dir(loc), filesystem.DirPermissionsDefault)

	return loc
}

// GetHome returns the home directory of the current user.
func (*Core) GetHome() string {
	home, _ := os.UserHomeDir()

	return home
}

// GetCacheRoot returns the cache root directory for the application.
func (obj *Core) GetCacheRoot() string {
	base, _ := os.UserCacheDir()

	loc := path.Join(base, obj.location[0])

	// XXX ignore errors?
	_ = os.MkdirAll(path.Dir(loc), filesystem.DirPermissionsDefault)

	return loc
}

// GetLogRoot returns the log root directory for the application.
func (obj *Core) GetLogRoot() string {
	var loc string

	base, _ := os.UserHomeDir()

	switch runtime.GOOS {
	case "darwin":
		loc = path.Join(base, "Library", "Logs", obj.location[0])
	default:
		loc = "/var/log/" + obj.location[0]
	}

	// XXX ignore errors?
	_ = os.MkdirAll(path.Dir(loc), filesystem.DirPermissionsDefault)

	return loc
}

func absolute(location ...string) string {
	loc := path.Join(location...)

	if !filepath.IsAbs(loc) {
		dir, _ := os.UserConfigDir()
		loc = path.Join(dir, loc)
	}

	return loc
}
