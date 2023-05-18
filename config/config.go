package config

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"

	"go.codecomet.dev/core/ca"
	"go.codecomet.dev/core/filesystem"
	"go.codecomet.dev/core/log"
	"go.codecomet.dev/core/network"
	"go.codecomet.dev/core/reporter"
	"go.codecomet.dev/core/telemetry"
)

func New(trustCA bool, appName string, location ...string) *Core {
	// Init filesystem first (capture the current, actual umask before we do anything)
	filesystem.Init()

	conf := &Core{
		Location: append([]string{appName}, location...),

		Client: &network.Config{
			TLSMin:              defaultTLSClientMinVersion,
			TLSHandshakeTimeout: defaultTLSHandshakeTimeout,
			DialerKeepAlive:     defaultDialerKeepAlive,
			DialerTimeout:       defaultDialerTimeout,
			DisallowSystemRoot:  false,
			CertPath:            defaultCertPath,
			KeyPath:             defaultKeyPath,
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

	if trustCA {
		conf.Server.ClientCA = ca.CodeComet
		conf.Client.RootCAs = []string{ca.CodeComet}
	}

	return conf
}

type CoreConfig interface {
	Resolve(...string) string
}

type Core struct {
	Reporter  *reporter.Config  `json:"reporter,omitempty"`
	Logger    *log.Config       `json:"logger,omitempty"`
	Telemetry *telemetry.Config `json:"telemetry,omitempty"`
	Client    *network.Config   `json:"client,omitempty"`
	Server    *network.Config   `json:"server,omitempty"`

	Umask int `json:"umask,omitempty"`

	Location []string `json:"-"`
}

func (obj *Core) Resolve(location ...string) string {
	// Get the absolute path of the containing dir of the config file, resolved against UserConfigDir
	base := absolute(obj.Location[:len(obj.Location)-1]...)

	loc := path.Join(location...)
	if !filepath.IsAbs(loc) {
		loc = path.Join(append([]string{base}, location...)...)
	}

	// XXX ignore errors?
	_ = os.MkdirAll(path.Dir(loc), filesystem.DirPermissionsDefault)

	return loc
}

func (obj *Core) Exist() bool {
	_, err := os.Stat(obj.Resolve(obj.Location...))

	return err == nil || !errors.Is(err, os.ErrNotExist)
}

func (obj *Core) Load(overload ...interface{}) error {
	var err error
	if len(overload) > 0 {
		err = read(overload[0], obj.Location...)

		field := reflect.ValueOf(overload[0]).Elem().FieldByName("Core")
		if field != (reflect.Value{}) {
			embed, ok := field.Interface().(*Core)
			if ok {
				filesystem.SetUmask(embed.Umask)
			}
		}

		return err
	}

	err = read(obj, obj.Location...)
	filesystem.SetUmask(obj.Umask)

	return err
}

func (obj *Core) Save(overload ...interface{}) error {
	if len(overload) > 0 {
		field := reflect.ValueOf(overload[0]).Elem().FieldByName("Core")
		if field != (reflect.Value{}) {
			embed, ok := field.Interface().(*Core)
			if ok {
				filesystem.SetUmask(embed.Umask)
			}
		}

		return write(overload[0], obj.Location...)
	}

	filesystem.SetUmask(obj.Umask)

	return write(obj, obj.Location...)
}

func (obj *Core) Remove() error {
	return remove(obj.Location...)
}

// XXX replace this with GetDataDir or GetCacheDir
/*
func (obj *Core) GetRunRoot() string {
	home, _ := os.UserHomeDir()
	loc := path.Join(home, "."+obj.Location[0], "run")

	// XXX ignore errors?
	_ = os.MkdirAll(path.Dir(loc), defaultDirPerms)

	return loc
}
*/

func (obj *Core) GetDataRoot() string {
	var loc string

	base, _ := os.UserHomeDir()

	switch runtime.GOOS {
	case "darwin":
		// XXX figure out impact on iCloud auto backup thing and containers
		loc = path.Join(base, "Library", "Application Support", obj.Location[0])
	default:
		loc = path.Join(base, "."+obj.Location[0])
	}

	// XXX ignore errors?
	_ = os.MkdirAll(path.Dir(loc), filesystem.DirPermissionsDefault)

	return loc
}

func (obj *Core) GetHome() string {
	home, _ := os.UserHomeDir()

	return home
}

func (obj *Core) GetCacheRoot() string {
	base, _ := os.UserCacheDir()

	loc := path.Join(base, obj.Location[0])

	// XXX ignore errors?
	_ = os.MkdirAll(path.Dir(loc), filesystem.DirPermissionsDefault)

	return loc
}

func (obj *Core) GetLogRoot() string {
	var loc string

	base, _ := os.UserHomeDir()

	switch runtime.GOOS {
	case "darwin":
		loc = path.Join(base, "Library", "Logs", obj.Location[0])
	default:
		loc = "/var/log/" + obj.Location[0]
	}

	// XXX ignore errors?
	_ = os.MkdirAll(path.Dir(loc), filesystem.DirPermissionsDefault)

	return loc
}
