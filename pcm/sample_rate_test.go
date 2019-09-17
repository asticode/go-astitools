package astipcm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSampleRateConverter(t *testing.T) {
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
	c := NewSampleRateConverter(1, 1, 1, sampleFunc)
	for _, s := range i {
		c.Add(s)
	}
	assert.Equal(t, i, o)

	// Simple src sample rate > dst sample rate
	o = []int{}
	c = NewSampleRateConverter(5, 3, 1, sampleFunc)
	for _, s := range i {
		c.Add(s)
	}
	assert.Equal(t, []int{1, 2, 4, 6, 7, 9, 11, 12, 14, 16, 17, 19}, o)

	// Multi channels
	o = []int{}
	c = NewSampleRateConverter(4, 2, 2, sampleFunc)
	for _, s := range i {
		c.Add(s)
	}
	assert.Equal(t, []int{1, 2, 4, 5, 8, 9, 12, 13, 16, 17}, o)

	// Realistic src sample rate > dst sample rate
	i = []int{}
	for idx := 0; idx < 4*44100; idx++ {
		i = append(i, idx+1)
	}
	o = []int{}
	c = NewSampleRateConverter(44100, 16000, 2, sampleFunc)
	for _, s := range i {
		c.Add(s)
	}
	assert.Len(t, o, 4*16000)

	// Create input
	i = []int{}
	for idx := 0; idx < 10; idx++ {
		i = append(i, idx+1)
	}

	// Simple src sample rate < dst sample rate
	o = []int{}
	c = NewSampleRateConverter(3, 5, 1, sampleFunc)
	for _, s := range i {
		c.Add(s)
	}
	assert.Equal(t, []int{1, 1, 2, 2, 3, 4, 4, 5, 5, 6, 7, 7, 8, 8, 9, 10, 10}, o)

	// Multi channels
	o = []int{}
	c = NewSampleRateConverter(3, 5, 2, sampleFunc)
	for _, s := range i {
		c.Add(s)
	}
	assert.Equal(t, []int{1, 2, 1, 2, 3, 4, 3, 4, 5, 6, 7, 8, 7, 8, 9, 10, 9, 10}, o)
}
