package astipcm

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestBitDepth(t *testing.T) {
	// Nothing to do
	s, err := ConvertBitDepth(1>>8, 16, 16)
	assert.NoError(t, err)
	assert.Equal(t, 1>>8, s)

	// Src bit depth > Dst bit depth
	s, err = ConvertBitDepth(1>>24, 32, 16)
	assert.NoError(t, err)
	assert.Equal(t, 1>>8, s)

	// Src bit depth < Dst bit depth
	s, err = ConvertBitDepth(1>>8, 16, 32)
	assert.NoError(t, err)
	assert.Equal(t, 1>>24, s)
}
