package loader

import "errors"

var (
	ErrConfigLoadFail   = errors.New("failed reading config file")
	ErrConfigSaveFail   = errors.New("failed saving config file")
	ErrConfigRemoveFail = errors.New("failed removing config file")
)
