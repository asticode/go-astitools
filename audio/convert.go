package astiaudio

import (
	"math"

	"github.com/pkg/errors"
)

// ConvertSampleRate converts the sample rate
func ConvertSampleRate(srcSamples []int32, srcSampleRate, dstSampleRate int) (dstSamples []int32, err error) {
	// Do nothing
	if srcSampleRate == dstSampleRate {
		dstSamples = srcSamples
		return
	}

	// For now we don't handle data loss
	if srcSampleRate < dstSampleRate {
		err = errors.New("astiaudio: src sample rate < dst sample rate")
		return
	}

	// Get ratio
	ratio := float64(srcSampleRate) / float64(dstSampleRate)

	// Loop
	var idx float64
	for {
		// Take floor of index
		floor := int(math.Floor(idx))

		// Out of range
		if floor >= len(srcSamples) {
			break
		}

		// Append sample
		dstSamples = append(dstSamples, srcSamples[floor])

		// Increment index
		idx += ratio
	}
	return
}

// ConvertBitDepth converts the bit depth
func ConvertBitDepth(srcSample int32, srcBitDepth, dstBitDepth int) (dstSample int32, err error) {
	// Nothing to do
	if srcBitDepth == dstBitDepth {
		return
	}

	// For now we don't handle data loss
	if srcBitDepth < dstBitDepth {
		err = errors.New("astiaudio: src bit depth < dst bit depth")
		return
	}

	// Convert
	dstSample = srcSample >> uint(srcBitDepth-dstBitDepth)
	return
}
