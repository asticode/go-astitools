package astiimage

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRGBA(t *testing.T) {
	c := &RGBA{}
	err := c.UnmarshalText([]byte("12345678"))
	assert.NoError(t, err)
	assert.Equal(t, RGBA{RGBA: color.RGBA{R: 0x78, G: 0x56, B: 0x34, A: 0x12}}, *c)
	err = c.UnmarshalText([]byte("9ABCDEFF"))
	assert.NoError(t, err)
	assert.Equal(t, RGBA{RGBA: color.RGBA{R: 0xff, G: 0xde, B: 0xbc, A: 0x9a}}, *c)
}
