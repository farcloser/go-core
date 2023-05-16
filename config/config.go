package config

import (
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"

	"go.codecomet.dev/core/ca"
	"go.codecomet.dev/core/log"
	"go.codecomet.dev/core/network"
	"go.codecomet.dev/core/reporter"
	"go.codecomet.dev/core/telemetry"
)

type Core struct {
	Location []string

	Umask     int               `json:"umask,omitempty"`
	Reporter  *reporter.Config  `json:"reporter,omitempty"`
	Logger    *log.Config       `json:"logger,omitempty"`
	Telemetry *telemetry.Config `json:"telemetry,omitempty"`
	Client    *network.Config   `json:"client,omitempty"`
	Server    *network.Config   `json:"server,omitempty"`
}

func New(trustCA bool, appName string, location ...string) *Core {
	cnf := &Core{
		Location: append([]string{appName}, location...),

		Umask: defaultUmask,

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

	umask(cnf.Umask)
	cnf.Client.Resolve = cnf.Resolve
	cnf.Server.Resolve = cnf.Resolve

	if trustCA {
		cnf.Server.ClientCA = ca.CodeComet
		cnf.Client.RootCAs = []string{ca.CodeComet}
	}

	return cnf
}

func (obj *Core) Load(overload ...interface{}) error {
	var err error
	if len(overload) > 0 {
		err = Read(overload[0], obj.Location...)

		field := reflect.ValueOf(overload[0]).Elem().FieldByName("Core")
		if field != (reflect.Value{}) {
			embed, ok := field.Interface().(*Core)
			if ok {
				umask(embed.Umask)
			}
		}

		return err
	}

	err = Read(obj, obj.Location...)
	umask(obj.Umask)

	return err
}

func (obj *Core) Save(overload ...interface{}) error {
	if len(overload) > 0 {
		field := reflect.ValueOf(overload[0]).Elem().FieldByName("Core")
		if field != (reflect.Value{}) {
			embed, ok := field.Interface().(*Core)
			if ok {
				umask(embed.Umask)
			}
		}

		return Write(overload[0], obj.Location...)
	}

	umask(obj.Umask)

	return Write(obj, obj.Location...)
}

func (obj *Core) Delete() error {
	return Delete(obj.Location...)
}

func (obj *Core) Resolve(location ...string) string {
	// Get the absolute path of the containing dir of the config file, resolved against UserConfigDir
	base := Absolute(obj.Location[:len(obj.Location)-1]...)

	loc := path.Join(location...)
	if !filepath.IsAbs(loc) {
		loc = path.Join(append([]string{base}, location...)...)
	}

	// XXX ignore errors?
	_ = os.MkdirAll(path.Dir(loc), DefaultDirPerms)

	return loc
}

// XXX replace this with GetDataDir or GetCacheDir
/*
func (obj *Core) GetRunRoot() string {
	home, _ := os.UserHomeDir()
	loc := path.Join(home, "."+obj.Location[0], "run")

	// XXX ignore errors?
	_ = os.MkdirAll(path.Dir(loc), DefaultDirPerms)

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
	_ = os.MkdirAll(path.Dir(loc), DefaultDirPerms)

	return loc
}

func (obj *Core) GetCacheRoot() string {
	base, _ := os.UserCacheDir()

	loc := path.Join(base, obj.Location[0])

	// XXX ignore errors?
	_ = os.MkdirAll(path.Dir(loc), DefaultDirPerms)

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
	_ = os.MkdirAll(path.Dir(loc), DefaultDirPerms)

	return loc
}
