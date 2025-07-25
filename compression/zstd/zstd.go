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

package zstd

import (
	"io"

	cmp "github.com/klauspost/compress/zstd"
)

type (
	// EncoderLevel predefines encoder compression levels.
	EncoderLevel = cmp.EncoderLevel
	// Encoder provides encoding to Zstandard.
	Encoder = cmp.Encoder
	// EOption is an option for creating a encoder.
	EOption = cmp.EOption
)

// EncoderLevelFromZstd converts a zstd compression level to an EncoderLevel.
func EncoderLevelFromZstd(level int) EncoderLevel {
	return cmp.EncoderLevelFromZstd(level)
}

// NewWriter creates a new Zstandard encoder that writes to the provided io.Writer.
//
//nolint:wrapcheck
func NewWriter(w io.Writer, opts ...EOption) (*Encoder, error) {
	return cmp.NewWriter(w, opts...)
}

// WithEncoderLevel sets the compression level for the encoder.
func WithEncoderLevel(l EncoderLevel) EOption {
	return cmp.WithEncoderLevel(l)
}
