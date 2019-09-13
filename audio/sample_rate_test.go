package astiaudio

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSampleRateConverter(t *testing.T) {
	// Create input
	var i []int32
	for idx := int32(0); idx < 20; idx++ {
		i = append(i, idx+1)
	}

	// Nothing to do
	var o []int32
	c := NewSampleRateConverter(1, 1, func(s int32) { o = append(o, s) })
	for _, s := range i {
		c.Add(s)
	}
	assert.Equal(t, i, o)

	// Simple src sample rate > dst sample rate
	o = []int32{}
	c = NewSampleRateConverter(5, 3, func(s int32) { o = append(o, s) })
	for _, s := range i {
		c.Add(s)
	}
	assert.Equal(t, []int32{1, 2, 4, 6, 7, 9, 11, 12, 14, 16, 17, 19}, o)

	// Realistic src sample rate > dst sample rate
	i = []int32{}
	for idx := int32(0); idx < 4*44100; idx++ {
		i = append(i, idx+1)
	}
	o = []int32{}
	c = NewSampleRateConverter(44100, 16000, func(s int32) { o = append(o, s) })
	for _, s := range i {
		c.Add(s)
	}
	assert.Len(t, o, 4*16000)
}
