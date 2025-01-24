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

package units_test

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"

	"go.farcloser.world/core/units"
)

func ExampleBytesSize() {
	fmt.Println(units.BytesSize(1024))
	fmt.Println(units.BytesSize(1024 * 1024))
	fmt.Println(units.BytesSize(1048576))
	fmt.Println(units.BytesSize(2 * units.MiB))
	fmt.Println(units.BytesSize(3.42 * units.GiB))
	fmt.Println(units.BytesSize(5.372 * units.TiB))
	fmt.Println(units.BytesSize(2.22 * units.PiB))
	// Output:
	// 1KiB
	// 1MiB
	// 1MiB
	// 2MiB
	// 3.42GiB
	// 5.372TiB
	// 2.22PiB
}

func ExampleHumanSize() {
	fmt.Println(units.HumanSize(1000))
	fmt.Println(units.HumanSize(1024))
	fmt.Println(units.HumanSize(1000000))
	fmt.Println(units.HumanSize(1048576))
	fmt.Println(units.HumanSize(2 * units.MB))
	fmt.Println(units.HumanSize(float64(3.42 * units.GB)))
	fmt.Println(units.HumanSize(float64(5.372 * units.TB)))
	fmt.Println(units.HumanSize(float64(2.22 * units.PB)))
	// Output:
	// 1kB
	// 1.024kB
	// 1MB
	// 1.049MB
	// 2MB
	// 3.42GB
	// 5.372TB
	// 2.22PB
}

func ExampleFromHumanSize() {
	fmt.Println(units.FromHumanSize("32"))
	fmt.Println(units.FromHumanSize("32b"))
	fmt.Println(units.FromHumanSize("32B"))
	fmt.Println(units.FromHumanSize("32k"))
	fmt.Println(units.FromHumanSize("32K"))
	fmt.Println(units.FromHumanSize("32kb"))
	fmt.Println(units.FromHumanSize("32Kb"))
	fmt.Println(units.FromHumanSize("32Mb"))
	fmt.Println(units.FromHumanSize("32Gb"))
	fmt.Println(units.FromHumanSize("32Tb"))
	fmt.Println(units.FromHumanSize("32Pb"))
	// Output:
	// 32 <nil>
	// 32 <nil>
	// 32 <nil>
	// 32000 <nil>
	// 32000 <nil>
	// 32000 <nil>
	// 32000 <nil>
	// 32000000 <nil>
	// 32000000000 <nil>
	// 32000000000000 <nil>
	// 32000000000000000 <nil>
}

func ExampleRAMInBytes() {
	fmt.Println(units.RAMInBytes("32"))
	fmt.Println(units.RAMInBytes("32b"))
	fmt.Println(units.RAMInBytes("32B"))
	fmt.Println(units.RAMInBytes("32k"))
	fmt.Println(units.RAMInBytes("32K"))
	fmt.Println(units.RAMInBytes("32kb"))
	fmt.Println(units.RAMInBytes("32Kb"))
	fmt.Println(units.RAMInBytes("32Mb"))
	fmt.Println(units.RAMInBytes("32Gb"))
	fmt.Println(units.RAMInBytes("32Tb"))
	fmt.Println(units.RAMInBytes("32Pb"))
	fmt.Println(units.RAMInBytes("32PB"))
	fmt.Println(units.RAMInBytes("32P"))
	// Output:
	// 32 <nil>
	// 32 <nil>
	// 32 <nil>
	// 32768 <nil>
	// 32768 <nil>
	// 32768 <nil>
	// 32768 <nil>
	// 33554432 <nil>
	// 34359738368 <nil>
	// 35184372088832 <nil>
	// 36028797018963968 <nil>
	// 36028797018963968 <nil>
	// 36028797018963968 <nil>
}

func TestBytesSize(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "1KiB", units.BytesSize(1024))
	assert.Equal(t, "1MiB", units.BytesSize(1024*1024))
	assert.Equal(t, "1MiB", units.BytesSize(1048576))
	assert.Equal(t, "2MiB", units.BytesSize(2*units.MiB))
	assert.Equal(t, "3.42GiB", units.BytesSize(3.42*units.GiB))
	assert.Equal(t, "5.372TiB", units.BytesSize(5.372*units.TiB))
	assert.Equal(t, "2.22PiB", units.BytesSize(2.22*units.PiB))
	assert.Equal(t, "1.049e+06YiB", units.BytesSize(units.KiB*units.KiB*units.KiB*units.KiB*units.KiB*units.PiB))
}

func TestHumanSize(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "1kB", units.HumanSize(1000))
	assert.Equal(t, "1.024kB", units.HumanSize(1024))
	assert.Equal(t, "1MB", units.HumanSize(1000000))
	assert.Equal(t, "1.049MB", units.HumanSize(1048576))
	assert.Equal(t, "2MB", units.HumanSize(2*units.MB))
	assert.Equal(t, "3.42GB", units.HumanSize(float64(3.42*units.GB)))
	assert.Equal(t, "5.372TB", units.HumanSize(float64(5.372*units.TB)))
	assert.Equal(t, "2.22PB", units.HumanSize(float64(2.22*units.PB)))
	assert.Equal(t, "1e+04YB", units.HumanSize(float64(10000000000000*units.PB)))
}

func TestFromHumanSize(t *testing.T) {
	t.Parallel()

	res, err := units.FromHumanSize("0")
	assert.NilError(t, err)
	assert.Equal(t, int64(0), res)

	res, err = units.FromHumanSize("0b")
	assert.NilError(t, err)
	assert.Equal(t, int64(0), res)

	res, err = units.FromHumanSize("0B")
	assert.NilError(t, err)
	assert.Equal(t, int64(0), res)

	res, err = units.FromHumanSize("0 B")
	assert.NilError(t, err)
	assert.Equal(t, int64(0), res)

	res, err = units.FromHumanSize("32")
	assert.NilError(t, err)
	assert.Equal(t, int64(32), res)

	res, err = units.FromHumanSize("32b")
	assert.NilError(t, err)
	assert.Equal(t, int64(32), res)

	res, err = units.FromHumanSize("32B")
	assert.NilError(t, err)
	assert.Equal(t, int64(32), res)

	res, err = units.FromHumanSize("32k")
	assert.NilError(t, err)
	assert.Equal(t, int64(32*units.KB), res)

	res, err = units.FromHumanSize("32K")
	assert.NilError(t, err)
	assert.Equal(t, int64(32*units.KB), res)

	res, err = units.FromHumanSize("32kb")
	assert.NilError(t, err)
	assert.Equal(t, int64(32*units.KB), res)

	res, err = units.FromHumanSize("32Kb")
	assert.NilError(t, err)
	assert.Equal(t, int64(32*units.KB), res)

	res, err = units.FromHumanSize("32Mb")
	assert.NilError(t, err)
	assert.Equal(t, int64(32*units.MB), res)

	res, err = units.FromHumanSize("32Gb")
	assert.NilError(t, err)
	assert.Equal(t, int64(32*units.GB), res)

	res, err = units.FromHumanSize("32Tb")
	assert.NilError(t, err)
	assert.Equal(t, int64(32*units.TB), res)

	res, err = units.FromHumanSize("32Pb")
	assert.NilError(t, err)
	assert.Equal(t, int64(32*units.PB), res)

	res, err = units.FromHumanSize("32.5kB")
	assert.NilError(t, err)
	assert.Equal(t, int64(32.5*units.KB), res)

	res, err = units.FromHumanSize("32.5 kB")
	assert.NilError(t, err)
	assert.Equal(t, int64(32.5*units.KB), res)

	res, err = units.FromHumanSize("32.5 B")
	assert.NilError(t, err)
	assert.Equal(t, int64(32), res)

	res, err = units.FromHumanSize("0.3 K")
	assert.NilError(t, err)
	assert.Equal(t, int64(300), res)

	res, err = units.FromHumanSize(".3kB")
	assert.NilError(t, err)
	assert.Equal(t, int64(300), res)

	res, err = units.FromHumanSize("0.")
	assert.NilError(t, err)
	assert.Equal(t, int64(0), res)

	res, err = units.FromHumanSize("0. ")
	assert.NilError(t, err)
	assert.Equal(t, int64(0), res)

	res, err = units.FromHumanSize("0.b")
	assert.NilError(t, err)
	assert.Equal(t, int64(0), res)

	res, err = units.FromHumanSize("0.B")
	assert.NilError(t, err)
	assert.Equal(t, int64(0), res)

	res, err = units.FromHumanSize("-0")
	assert.NilError(t, err)
	assert.Equal(t, int64(0), res)

	res, err = units.FromHumanSize("-0b")
	assert.NilError(t, err)
	assert.Equal(t, int64(0), res)

	res, err = units.FromHumanSize("-0B")
	assert.NilError(t, err)
	assert.Equal(t, int64(0), res)

	res, err = units.FromHumanSize("-0 b")
	assert.NilError(t, err)
	assert.Equal(t, int64(0), res)

	res, err = units.FromHumanSize("-0 B")
	assert.NilError(t, err)
	assert.Equal(t, int64(0), res)

	res, err = units.FromHumanSize("32.")
	assert.NilError(t, err)
	assert.Equal(t, int64(32), res)

	res, err = units.FromHumanSize("32.b")
	assert.NilError(t, err)
	assert.Equal(t, int64(32), res)

	res, err = units.FromHumanSize("32.B")
	assert.NilError(t, err)
	assert.Equal(t, int64(32), res)

	res, err = units.FromHumanSize("32. b")
	assert.NilError(t, err)
	assert.Equal(t, int64(32), res)

	res, err = units.FromHumanSize("32. B")
	assert.NilError(t, err)
	assert.Equal(t, int64(32), res)

	// We do not tolerate extra leading or trailing spaces
	// (except for a space after the number and a missing suffix).
	res, err = units.FromHumanSize("0 ")
	assert.NilError(t, err)
	assert.Equal(t, int64(0), res)

	_, err = units.FromHumanSize(" 0")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize(" 0b")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize(" 0B")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize(" 0 B")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize("0b ")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize("0B ")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize("0 B ")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize("")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize("hello")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize(".")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize(". ")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize(" ")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize("  ")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize(" .")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize(" . ")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize("-32")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize("-32b")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize("-32B")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize("-32 b")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize("-32 B")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize("32b.")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize("32B.")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize("32 b.")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize("32 B.")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize("32 bb")
	assert.ErrorIs(t, err, units.ErrInvalidSuffix)

	_, err = units.FromHumanSize("32 BB")
	assert.ErrorIs(t, err, units.ErrInvalidSuffix)

	_, err = units.FromHumanSize("32 b b")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize("32 B B")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize("32  b")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize("32  B")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize(" 32 ")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize("32m b")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.FromHumanSize("32bm")
	assert.ErrorIs(t, err, units.ErrInvalidSuffix)
}

func TestRAMInBytes(t *testing.T) {
	t.Parallel()

	res, err := units.RAMInBytes("32")
	assert.NilError(t, err)
	assert.Equal(t, int64(32), res)

	res, err = units.RAMInBytes("32b")
	assert.NilError(t, err)
	assert.Equal(t, int64(32), res)

	res, err = units.RAMInBytes("32B")
	assert.NilError(t, err)
	assert.Equal(t, int64(32), res)

	res, err = units.RAMInBytes("32k")
	assert.NilError(t, err)
	assert.Equal(t, int64(32*units.KiB), res)

	res, err = units.RAMInBytes("32K")
	assert.NilError(t, err)
	assert.Equal(t, int64(32*units.KiB), res)

	res, err = units.RAMInBytes("32kb")
	assert.NilError(t, err)
	assert.Equal(t, int64(32*units.KiB), res)

	res, err = units.RAMInBytes("32Kb")
	assert.NilError(t, err)
	assert.Equal(t, int64(32*units.KiB), res)

	res, err = units.RAMInBytes("32Kib")
	assert.NilError(t, err)
	assert.Equal(t, int64(32*units.KiB), res)

	res, err = units.RAMInBytes("32KIB")
	assert.NilError(t, err)
	assert.Equal(t, int64(32*units.KiB), res)

	res, err = units.RAMInBytes("32Mb")
	assert.NilError(t, err)
	assert.Equal(t, int64(32*units.MiB), res)

	res, err = units.RAMInBytes("32Gb")
	assert.NilError(t, err)
	assert.Equal(t, int64(32*units.GiB), res)

	res, err = units.RAMInBytes("32Tb")
	assert.NilError(t, err)
	assert.Equal(t, int64(32*units.TiB), res)

	res, err = units.RAMInBytes("32Pb")
	assert.NilError(t, err)
	assert.Equal(t, int64(32*units.PiB), res)

	res, err = units.RAMInBytes("32PB")
	assert.NilError(t, err)
	assert.Equal(t, int64(32*units.PiB), res)

	res, err = units.RAMInBytes("32P")
	assert.NilError(t, err)
	assert.Equal(t, int64(32*units.PiB), res)

	res, err = units.RAMInBytes("32.3")
	assert.NilError(t, err)
	assert.Equal(t, int64(32), res)

	tmp := 32.3 * units.MiB
	res, err = units.RAMInBytes("32.3 mb")
	assert.NilError(t, err)
	assert.Equal(t, int64(tmp), res)

	tmp = 0.3 * units.MiB
	res, err = units.RAMInBytes("0.3MB")
	assert.NilError(t, err)
	assert.Equal(t, int64(tmp), res)

	_, err = units.RAMInBytes("")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.RAMInBytes("hello")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.RAMInBytes("-32")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.RAMInBytes(" 32 ")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.RAMInBytes("32m b")
	assert.ErrorIs(t, err, units.ErrInvalidSize)

	_, err = units.RAMInBytes("32bm")
	assert.ErrorIs(t, err, units.ErrInvalidSuffix)
}

func BenchmarkParseSize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, s := range []string{
			"", "32", "32b", "32 B", "32k", "32.5 K", "32kb", "32 Kb",
			"32.8Mb", "32.9Gb", "32.777Tb", "32Pb", "0.3Mb", "-1",
		} {
			_, _ = units.FromHumanSize(s)
			_, _ = units.RAMInBytes(s)
		}
	}
}
