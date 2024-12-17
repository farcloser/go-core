package filesystem_test

import (
	"fmt"
	"runtime"
	"testing"

	"gotest.tools/v3/assert"

	"go.farcloser.world/core/filesystem"
)

func TestFilesystemRestrictions(t *testing.T) {
	t.Parallel()

	invalid := []string{
		"/",
		"/start",
		"mid/dle",
		"end/",
		".",
		"..",
		"",
		fmt.Sprintf("A%0255s", "A"),
	}

	valid := []string{
		fmt.Sprintf("A%0254s", "A"),
		"test",
		"test-hyphen",
		".start.dot",
		"mid.dot",
		"∞",
	}

	if runtime.GOOS == "windows" {
		invalid = append(invalid, []string{
			"\\start",
			"mid\\dle",
			"end\\",
			"\\",
			"\\.",
			"com².whatever",
			"lpT2",
			"Prn.",
			"nUl",
			"AUX",
			"A<A",
			"A>A",
			"A:A",
			"A\"A",
			"A|A",
			"A?A",
			"A*A",
			"end.dot.",
			"end.space ",
		}...)
	}

	for _, v := range invalid {
		err := filesystem.ValidatePathComponent(v)
		assert.ErrorIs(t, err, filesystem.ErrInvalidPath, v)
	}

	for _, v := range valid {
		err := filesystem.ValidatePathComponent(v)
		assert.NilError(t, err, v)
	}
}
