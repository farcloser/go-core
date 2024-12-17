package filesystem

import (
	"errors"
	"strings"
)

var (
	errForbiddenChars    = errors.New("forbidden characters in path component")
	errForbiddenKeywords = errors.New("forbidden keywords in path component")
)

// ValidatePathComponent will enforce os specific filename restrictions on a single path component.
func ValidatePathComponent(pathComponent string) error {
	// https://en.wikipedia.org/wiki/Comparison_of_file_systems#Limits
	if len(pathComponent) > pathComponentMaxLength {
		return errors.Join(ErrInvalidPath, errInvalidPathTooLong)
	}

	if strings.TrimSpace(pathComponent) == "" {
		return errors.Join(ErrInvalidPath, errInvalidPathEmpty)
	}

	if err := validatePlatformSpecific(pathComponent); err != nil {
		return errors.Join(ErrInvalidPath, err)
	}

	return nil
}
