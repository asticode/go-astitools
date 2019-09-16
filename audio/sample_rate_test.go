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

	// Create sample func
	var o []int32
	var sampleFunc = func(s int32) (err error) {
		o = append(o, s)
		return
	}

	// Nothing to do
	c := NewSampleRateConverter(1, 1, sampleFunc)
	for _, s := range i {
		c.Add(s)
	}
	assert.Equal(t, i, o)

	// Simple src sample rate > dst sample rate
	o = []int32{}
	c = NewSampleRateConverter(5, 3, sampleFunc)
	for _, s := range i {
		c.Add(s)
	}
	assert.Equal(t, []int32{1, 2, 3, 5, 7, 8, 10, 12, 13, 15, 17, 18}, o)

	// Realistic src sample rate > dst sample rate
	i = []int32{}
	for idx := int32(0); idx < 4*44100; idx++ {
		i = append(i, idx+1)
	}
	o = []int32{}
	c = NewSampleRateConverter(44100, 16000, sampleFunc)
	for _, s := range i {
		c.Add(s)
	}
	assert.Len(t, o, 4*16000)
}
