package loader

import (
	"errors"
	"os"
)

// Inheriting CoreConfig gives default implementation
type IConfiguration interface {
	/*
		// Resolve returns a resolved path relative to the config file location
		Resolve(all ...string) string
		// Ensure makes sure the parent directory of the resolved path exists
		Ensure(all ...string) error

		// GetHome returns the user home directory
		GetHome() string
		// GetDataRoot returns the app persistent storage location
		GetDataRoot() string
		// GetCacheRoot returns the app transient storage location
		GetCacheRoot() string
		// GetLogRoot returns the app logs location
		GetLogRoot() string
	*/

	// OnIO is where you should implement logic to be done after loading and before saving
	// If you inherit CoreConfig, be sure to call super.OnIO() first
	OnIO()
	// GetLocation returns the relative location of the config file (still has to be resolved for a full path)
	GetLocation() []string
}

func Exist(obj IConfiguration) bool {
	_, err := os.Stat(absolute(obj.GetLocation()...))

	return err == nil || !errors.Is(err, os.ErrNotExist)
}

func Load(obj IConfiguration) error {
	err := read(obj, obj.GetLocation()...)
	if err != nil {
		return err
	}

	obj.OnIO()

	return nil
}

func Save(obj IConfiguration) error {
	obj.OnIO()

	return write(obj, obj.GetLocation()...)
}

func Remove(obj IConfiguration) error {
	return remove(obj.GetLocation()...)
}
