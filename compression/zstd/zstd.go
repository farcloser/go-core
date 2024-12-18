package zstd

import (
	"io"

	cmp "github.com/klauspost/compress/zstd"
)

type (
	EncoderLevel = cmp.EncoderLevel
	Encoder      = cmp.Encoder
	EOption      = cmp.EOption
)

func EncoderLevelFromZstd(level int) EncoderLevel {
	return cmp.EncoderLevelFromZstd(level)
}

func NewWriter(w io.Writer, opts ...EOption) (*Encoder, error) {
	return cmp.NewWriter(w, opts...)
}

func WithEncoderLevel(l EncoderLevel) EOption {
	return cmp.WithEncoderLevel(l)
}
