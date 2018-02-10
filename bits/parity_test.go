package astibits

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testParity(i uint8) bool {
	return (i&0x1)^(i>>1&0x1)^(i>>2&0x1)^(i>>3&0x1)^(i>>4&0x1)^(i>>5&0x1)^(i>>6&0x1)^(i>>7&0x1) > 0
}

func TestParity(t *testing.T) {
	for i := 0; i < 256; i++ {
		v, okV := Parity(uint8(i))
		okE := testParity(uint8(i))
		if !okE {
			assert.False(t, okV)
		} else {
			assert.True(t, okV)
			assert.Equal(t, uint8(i)&0x7f, v)
		}
	}
}
