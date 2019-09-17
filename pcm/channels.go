package astipcm

import (
	"github.com/pkg/errors"
)

type ChannelsConverter struct {
	dstNumChannels int
	fn             SampleFunc
	srcNumChannels int
	srcSamples     int
}

func NewChannelsConverter(srcNumChannels, dstNumChannels int, fn SampleFunc) *ChannelsConverter {
	return &ChannelsConverter{
		dstNumChannels: dstNumChannels,
		fn:             fn,
		srcNumChannels: srcNumChannels,
	}
}

func (c *ChannelsConverter) Reset() {
	c.srcSamples = 0
}

func (c *ChannelsConverter) Add(i int) (err error) {
	// Forward sample
	if c.srcNumChannels == c.dstNumChannels {
		if err = c.fn(i); err != nil {
			err = errors.Wrap(err, "astipcm: handling sample failed")
			return
		}
		return
	}

	// Reset
	if c.srcSamples == c.srcNumChannels {
		c.srcSamples = 0
	}

	// Increment src samples
	c.srcSamples++

	// Throw away data
	if c.srcNumChannels > c.dstNumChannels {
		// Throw away sample
		if c.srcSamples > c.dstNumChannels {
			return
		}

		// Custom
		if err = c.fn(i); err != nil {
			err = errors.Wrap(err, "astipcm: handling sample failed")
			return
		}
		return
	}

	// Store
	var ss []int
	if c.srcSamples < c.srcNumChannels {
		ss = []int{i}
	} else {
		// Repeat data
		for idx := c.srcNumChannels; idx <= c.dstNumChannels; idx++ {
			ss = append(ss, i)
		}
	}

	// Loop through samples
	for _, s := range ss {
		// Custom
		if err = c.fn(s); err != nil {
			err = errors.Wrap(err, "astipcm: handling sample failed")
			return
		}
	}
	return
}
