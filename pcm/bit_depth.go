package astipcm

// ConvertBitDepth converts the bit depth
func ConvertBitDepth(srcSample int, srcBitDepth, dstBitDepth int) (dstSample int, err error) {
	// Nothing to do
	if srcBitDepth == dstBitDepth {
		dstSample = srcSample
		return
	}

	// Convert
	if srcBitDepth < dstBitDepth {
		dstSample = srcSample << uint(dstBitDepth-srcBitDepth)
	} else {
		dstSample = srcSample >> uint(srcBitDepth-dstBitDepth)
	}
	return
}
