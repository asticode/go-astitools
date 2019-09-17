package astipcm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChannelsConverter(t *testing.T) {
	// Create input
	var i []int
	for idx := 0; idx < 20; idx++ {
		i = append(i, idx+1)
	}

	// Create sample func
	var o []int
	var sampleFunc = func(s int) (err error) {
		o = append(o, s)
		return
	}

	// Nothing to do
	c := NewChannelsConverter(3, 3, sampleFunc)
	for _, s := range i {
		c.Add(s)
	}
	assert.Equal(t, i, o)

	// Throw away data
	o = []int{}
	c = NewChannelsConverter(3, 1, sampleFunc)
	for _, s := range i {
		c.Add(s)
	}
	assert.Equal(t, []int{1, 4, 7, 10, 13, 16, 19}, o)

	// Repeat data
	o = []int{}
	c = NewChannelsConverter(1, 2, sampleFunc)
	for _, s := range i {
		c.Add(s)
	}
	assert.Equal(t, []int{1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7, 8, 8, 9, 9, 10, 10, 11, 11, 12, 12, 13, 13, 14, 14, 15, 15, 16, 16, 17, 17, 18, 18, 19, 19, 20, 20}, o)
}
