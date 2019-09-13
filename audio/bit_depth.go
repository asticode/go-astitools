package astiaudio

import (
	"github.com/pkg/errors"
)

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
