package astipcm

import (
	"fmt"

	"github.com/pkg/errors"
)

type SampleFunc func(s int) error

type SampleRateConverter struct {
	b                    [][]int
	dstSampleRate        int
	fn                   SampleFunc
	numChannels          int
	numChannelsProcessed int
	numSamplesOutputed   int
	numSamplesProcessed  int
	srcSampleRate        int
}

func NewSampleRateConverter(srcSampleRate, dstSampleRate, numChannels int, fn SampleFunc) *SampleRateConverter {
	return &SampleRateConverter{
		b:             make([][]int, numChannels),
		dstSampleRate: dstSampleRate,
		fn:            fn,
		numChannels:   numChannels,
		srcSampleRate: srcSampleRate,
	}
}

func (c *SampleRateConverter) Reset() {
	c.b = make([][]int, c.numChannels)
	c.numChannelsProcessed = 0
	c.numSamplesOutputed = 0
	c.numSamplesProcessed = 0
}

func (c *SampleRateConverter) Add(i int) (err error) {
	// Forward sample
	if c.srcSampleRate == c.dstSampleRate {
		if err = c.fn(i); err != nil {
			err = errors.Wrap(err, "astipcm: handling sample failed")
			return
		}
		return
	}

	// Increment num channels processed
	c.numChannelsProcessed++

	// Reset num channels processed
	if c.numChannelsProcessed > c.numChannels {
		c.numChannelsProcessed = 1
	}

	// Only increment num samples processed if all channels have been processed
	if c.numChannelsProcessed == c.numChannels {
		c.numSamplesProcessed++
	}

	// Append sample to buffer
	c.b[c.numChannelsProcessed-1] = append(c.b[c.numChannelsProcessed-1], i)

	// Throw away data
	if c.srcSampleRate > c.dstSampleRate {
		// Make sure to always keep the first sample but do nothing until we have all channels or target sample has been
		// reached
		if (c.numSamplesOutputed > 0 && float64(c.numSamplesProcessed) < 1.0+float64(c.numSamplesOutputed)*float64(c.srcSampleRate)/float64(c.dstSampleRate)) || c.numChannelsProcessed < c.numChannels {
			return
		}

		// Loop through channels
		for idx, b := range c.b {
			// Merge samples
			var s int
			for _, v := range b {
				s += v
			}
			s /= len(b)

			// Reset buffer
			c.b[idx] = []int{}

			// Custom
			if err = c.fn(s); err != nil {
				err = errors.Wrap(err, "astipcm: handling sample failed")
				return
			}
		}

		// Increment num samples outputed
		c.numSamplesOutputed++
		return
	}

	// Do nothing until we have all channels
	if c.numChannelsProcessed < c.numChannels {
		return
	}

	// Repeat data
	for c.numSamplesOutputed == 0 || float64(c.numSamplesProcessed)+1.0 > 1.0+float64(c.numSamplesOutputed)*float64(c.srcSampleRate)/float64(c.dstSampleRate) {
		// Loop through channels
		for _, b := range c.b {
			// Invalid length
			if len(b) != 1 {
				err = fmt.Errorf("astipcm: invalid buffer item length %d", len(b))
				return
			}

			// Custom
			if err = c.fn(b[0]); err != nil {
				err = errors.Wrap(err, "astipcm: handling sample failed")
				return
			}
		}

		// Increment num samples outputed
		c.numSamplesOutputed++
	}

	// Reset buffer
	c.b = make([][]int, c.numChannels)
	return
}
