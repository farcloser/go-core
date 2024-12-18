package semver

import (
	"errors"

	"github.com/Masterminds/semver/v3"
)

type (
	Version     = semver.Version
	Constraints = semver.Constraints
)

var (
	// ErrInvalidSemVer is returned a version is found to be invalid when
	// being parsed.
	ErrInvalidSemVer = semver.ErrInvalidSemVer

	// ErrSegmentStartsZero is returned when a version segment starts with 0.
	// This is invalid in SemVer.
	ErrSegmentStartsZero = semver.ErrSegmentStartsZero

	// ErrInvalidMetadata is returned when the metadata is an invalid format.
	ErrInvalidMetadata = semver.ErrInvalidMetadata

	// ErrInvalidPrerelease is returned when the pre-release is an invalid format.
	ErrInvalidPrerelease = semver.ErrInvalidPrerelease
)

func NewVersion(str string) (*Version, error) {
	// ErrInvalidSemVer, or fmt.Errorf("Error parsing version segment: %s", errFromStrConv)
	// ErrInvalidMetadata,
	// ErrSegmentStartsZero,
	// ErrInvalidPrerelease
	version, err := semver.NewVersion(str)
	if err != nil &&
		!errors.Is(err, ErrInvalidSemVer) && !errors.Is(err, ErrSegmentStartsZero) &&
		!errors.Is(err, ErrInvalidMetadata) && !errors.Is(err, ErrInvalidPrerelease) {
		err = ErrInvalidSemVer
	}

	return version, err
}

func NewConstraint(str string) (*Constraints, error) {
	return semver.NewConstraint(str)
}
