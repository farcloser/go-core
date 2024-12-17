//go:build darwin || dragonfly || freebsd || illumos || linux || netbsd || openbsd

package filesystem

import (
	"fmt"
	"regexp"
)

// Note that Darwin has different restrictions on colons.
// https://stackoverflow.com/questions/1976007/what-characters-are-forbidden-in-windows-and-linux-directory-names
var (
	disallowedKeywords = regexp.MustCompile(`^([.]|[.][.])$`)
	reservedCharacters = regexp.MustCompile(`[\x{0}/]`)
)

func validatePlatformSpecific(pathComponent string) error {
	if reservedCharacters.MatchString(pathComponent) {
		return fmt.Errorf("%w: %q (%q)", errForbiddenChars, pathComponent, reservedCharacters)
	}

	if disallowedKeywords.MatchString(pathComponent) {
		return fmt.Errorf("%w: %q (%q)", errForbiddenKeywords, pathComponent, disallowedKeywords)
	}

	return nil
}
