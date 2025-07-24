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

// Package units forked from https://github.com/docker/go-units under Apache License	2.0
//
//revive:disable:add-constant
package units

import (
	"fmt"
	"strconv"
	"strings"
)

// See: http://en.wikipedia.org/wiki/Binary_prefix
const (
	// Decimal.

	KB = 1000      //nolint:varnamelen
	MB = 1000 * KB //nolint:varnamelen
	GB = 1000 * MB //nolint:varnamelen
	TB = 1000 * GB //nolint:varnamelen
	PB = 1000 * TB //nolint:varnamelen

	// Binary.

	KiB = 1024
	MiB = 1024 * KiB
	GiB = 1024 * MiB
	TiB = 1024 * GiB
	PiB = 1024 * TiB

	suffixTooLong = 3
	precision     = 4
	base          = 1024.0
	decimalBase   = 1000
)

type unitMap map[byte]int64

//nolint:gochecknoglobals
var (
	decimalMap = unitMap{
		'k': KB,
		'm': MB,
		'g': GB,
		't': TB,
		'p': PB,
	}
	binaryMap = unitMap{
		'k': KiB,
		'm': MiB,
		'g': GiB,
		't': TiB,
		'p': PiB,
	}
	decimapAbbrs = []string{
		"B",
		"kB",
		"MB",
		"GB",
		"TB",
		"PB",
		"EB",
		"ZB",
		"YB",
	}
	binaryAbbrs = []string{
		"B",
		"KiB",
		"MiB",
		"GiB",
		"TiB",
		"PiB",
		"EiB",
		"ZiB",
		"YiB",
	}
)

func getSizeAndUnit(size, base float64, _map []string) (float64, string) {
	index := 0

	unitsLimit := len(_map) - 1
	for size >= base && index < unitsLimit {
		size /= base
		index++
	}

	return size, _map[index]
}

// CustomSize returns a human-readable approximation of a size
// using custom format.
func CustomSize(format string, size, base float64, _map []string) string {
	size, unit := getSizeAndUnit(size, base, _map)

	return fmt.Sprintf(format, size, unit)
}

// HumanSizeWithPrecision allows the size to be in any precision,
// instead of 4 digit precision used in units.HumanSize.
func HumanSizeWithPrecision(size float64, precision int) string {
	size, unit := getSizeAndUnit(size, decimalBase, decimapAbbrs)

	return fmt.Sprintf("%.*g%s", precision, size, unit)
}

// HumanSize returns a human-readable approximation of a size
// capped at 4 valid numbers (e.g. "2.746 MB", "796 KB").
func HumanSize(size float64) string {
	return HumanSizeWithPrecision(size, precision)
}

// BytesSize returns a human-readable size in bytes, kibibytes,
// mebibytes, gibibytes, or tebibytes (e.g. "44kiB", "17MiB").
func BytesSize(size float64) string {
	return CustomSize("%.4g%s", size, base, binaryAbbrs)
}

// FromHumanSize returns an integer from a human-readable specification of a
// size using SI standard (e.g. "44kB", "17MB").
func FromHumanSize(size string) (int64, error) {
	return parseSize(size, decimalMap)
}

// RAMInBytes parses a human-readable string representing an amount of RAM
// in bytes, kibibytes, mebibytes, gibibytes, or tebibytes and
// returns the number of bytes, or -1 if the string is unparseable.
// Units are case-insensitive, and the 'b' suffix is optional.
func RAMInBytes(size string) (int64, error) {
	return parseSize(size, binaryMap)
}

// Parses the human-readable size string into the amount it represents.
func parseSize(sizeStr string, uMap unitMap) (int64, error) {
	// TODO: rewrite to use strings.Cut if there's a space
	// once Go < 1.18 is deprecated.
	sep := strings.LastIndexAny(sizeStr, "01234567890. ")
	if sep == -1 {
		// There should be at least a digit.
		return -1, fmt.Errorf("%w: '%s'", ErrInvalidSize, sizeStr)
	}

	var num, sfx string
	if sizeStr[sep] != ' ' {
		num = sizeStr[:sep+1]
		sfx = sizeStr[sep+1:]
	} else {
		// Omit the space separator.
		num = sizeStr[:sep]
		sfx = sizeStr[sep+1:]
	}

	size, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return -1, fmt.Errorf("%w: '%w'", ErrInvalidSize, err)
	}
	// Backward compatibility: reject negative sizes.
	if size < 0 {
		return -1, fmt.Errorf("%w: '%s'", ErrInvalidSize, sizeStr)
	}

	if len(sfx) == 0 {
		return int64(size), nil
	}

	// Process the suffix.

	if len(sfx) > suffixTooLong { // Too long.
		return -1, fmt.Errorf("%w: '%s'", ErrInvalidSuffix, sfx)
	}

	sfx = strings.ToLower(sfx)
	// Trivial case: b suffix.
	if sfx[0] == 'b' {
		if len(sfx) > 1 { // no extra characters allowed after b.
			return -1, fmt.Errorf("%w: '%s'", ErrInvalidSuffix, sfx)
		}

		return int64(size), nil
	}

	// A suffix from the map.
	mul, ok := uMap[sfx[0]]
	if !ok {
		return -1, fmt.Errorf("%w: '%s'", ErrInvalidSuffix, sfx)
	}

	size *= float64(mul)

	// The suffix may have extra "b" or "ib" (e.g. KiB or MB).
	switch {
	case len(sfx) == 2 && sfx[1] != 'b':
		return -1, fmt.Errorf("%w: '%s'", ErrInvalidSuffix, sfx)
	case len(sfx) == 3 && sfx[1:] != "ib":
		return -1, fmt.Errorf("%w: '%s'", ErrInvalidSuffix, sfx)
	}

	return int64(size), nil
}
