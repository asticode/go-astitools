package astiimage

import (
	"image/color"
	"strconv"

	"github.com/pkg/errors"
)

// RGBA wraps color.RGBA
type RGBA struct {
	color.RGBA
}

// NewRGBA creates a new RGBA color
func NewRGBA(a, b, g, r uint8) *RGBA {
	return &RGBA{RGBA: color.RGBA{
		A: a,
		B: b,
		G: g,
		R: r,
	}}
}

// UnmarshalText implements the encoding.TextUnmarshaler interface
func (r *RGBA) UnmarshalText(i []byte) (err error) {
	var p uint64
	if p, err = strconv.ParseUint(string(i), 16, 32); err != nil {
		err = errors.Wrapf(err, "astiimage: parsing uint %s failed", i)
		return
	}
	r.R = uint8(p & 0xff)
	r.G = uint8(p >> 8 & 0xff)
	r.B = uint8(p >> 16 & 0xff)
	r.A = uint8(p >> 24 & 0xff)
	return
}
