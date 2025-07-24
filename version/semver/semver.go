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

package semver

import (
	"errors"

	"github.com/Masterminds/semver/v3"
)

type (
	// Version represents a single semantic version.
	Version = semver.Version
	// Constraints represents a set of semantic version constraints.
	Constraints = semver.Constraints
)

// ErrInvalidVersion is returned when a version string is invalid or cannot be parsed.
var ErrInvalidVersion = errors.New("invalid version")

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

// NewVersion creates a new semantic version from a string.
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

// NewConstraint creates a new semantic version constraint from a string.
func NewConstraint(str string) (*Constraints, error) {
	constraint, err := semver.NewConstraint(str)
	if err != nil {
		err = errors.Join(ErrInvalidVersion, err)
	}

	return constraint, err
}
