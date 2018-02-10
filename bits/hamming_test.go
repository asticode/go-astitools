package astibits

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testHamming84Decode(i uint8) (o uint8, err error) {
	p1, d1, p2, d2, p3, d3, p4, d4 := i>>7&0x1, i>>6&0x1, i>>5&0x1, i>>4&0x1, i>>3&0x1, i>>2&0x1, i>>1&0x1, i&0x1
	testA := p1^d1^d3^d4 > 0
	testB := d1^p2^d2^d4 > 0
	testC := d1^d2^p3^d3 > 0
	testD := p1^d1^p2^d2^p3^d3^p4^d4 > 0
	if testA && testB && testC {
		// p4 may be incorrect
	} else if testD && (!testA || !testB || !testC) {
		err = fmt.Errorf("hamming 8/4 decode of %.8b failed", i)
		return
	} else {
		if !testA && testB && testC {
			// p1 is incorrect
		} else if testA && !testB && testC {
			// p2 is incorrect
		} else if testA && testB && !testC {
			// p3 is incorrect
		} else if !testA && !testB && testC {
			// d4 is incorrect
			d4 ^= 1
		} else if testA && !testB && !testC {
			// d2 is incorrect
			d2 ^= 1
		} else if !testA && testB && !testC {
			// d3 is incorrect
			d3 ^= 1
		} else {
			// d1 is incorrect
			d1 ^= 1
		}
	}
	o = uint8(d4<<3 | d3<<2 | d2<<1 | d1)
	return
}

func TestHamming84Decode(t *testing.T) {
	for i := 0; i < 256; i++ {
		v, errV := Hamming84Decode(uint8(i))
		e, errE := testHamming84Decode(uint8(i))
		if errE != nil {
			assert.Error(t, errV)
		} else {
			assert.NoError(t, errV)
			assert.Equal(t, e, v)
		}
	}
}
