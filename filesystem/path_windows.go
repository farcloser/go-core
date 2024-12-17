package filesystem

import (
	"errors"
	"fmt"
	"regexp"
)

// See https://learn.microsoft.com/en-us/windows/win32/fileio/naming-a-file
// https://stackoverflow.com/questions/1976007/what-characters-are-forbidden-in-windows-and-linux-directory-names
var (
	disallowedKeywords = regexp.MustCompile(`(?i)^(con|prn|nul|aux|com[1-9¹²³]|lpt[1-9¹²³])([.].*)?$`)
	reservedCharacters = regexp.MustCompile(`[\x{0}-\x{1f}<>:"/\\|?*]`)

	errNoEndingSpaceDot = errors.New("component cannot end with a space or dot")
)

func validatePlatformSpecific(pathComponent string) error {
	if reservedCharacters.MatchString(pathComponent) {
		return fmt.Errorf("%w: %q (%q)", errForbiddenChars, pathComponent, reservedCharacters)
	}

	if disallowedKeywords.MatchString(pathComponent) {
		return fmt.Errorf("%w: %q (%q)", errForbiddenKeywords, pathComponent, disallowedKeywords)
	}

	if pathComponent[len(pathComponent)-1:] == "." || pathComponent[len(pathComponent)-1:] == " " {
		return fmt.Errorf("%w: %q", errNoEndingSpaceDot, pathComponent)
	}

	return nil
}
