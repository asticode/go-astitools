package astiaudio

import "github.com/pkg/errors"

type SampleFunc func(s int32) error

type SampleRateConverter struct {
	delta      float64
	dstSamples int
	fn         SampleFunc
	srcSamples int
}

func NewSampleRateConverter(srcSampleRate, dstSampleRate float64, fn SampleFunc) *SampleRateConverter {
	return &SampleRateConverter{
		delta: float64(srcSampleRate) / float64(dstSampleRate),
		fn:    fn,
	}
}

func (c *SampleRateConverter) Reset() {
	c.dstSamples = 0
	c.srcSamples = 0
}

func (c *SampleRateConverter) Add(i int32) (err error) {
	// Nothing to do
	if c.delta == 1 {
		c.fn(i)
		return
	}

	// Increment src samples
	c.srcSamples++

	// Throw away data
	if c.delta > 1 {
		// Make sure to always keep the first sample
		if c.dstSamples > 0 && float64(c.srcSamples) <= c.delta*float64(c.dstSamples) {
			return
		}

		// Increment dst samples
		c.dstSamples++

		// Custom
		if err = c.fn(i); err != nil {
			err = errors.Wrap(err, "astiaudio: handling sample failed")
			return
		}
		return
	}

	// TODO Repeat data
	panic("astiaudio: repeating data when converting sample rate is not implemented")
}
