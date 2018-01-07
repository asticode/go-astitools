package astiaudio

import "math"

// AudioLevel computes the audio level of samples
// https://dsp.stackexchange.com/questions/2951/loudness-of-pcm-stream
// https://dsp.stackexchange.com/questions/290/getting-loudness-of-a-track-with-rms?noredirect=1&lq=1
func AudioLevel(samples []int32) float64 {
	// Compute sum of square values
	var sum float64
	for _, s := range samples {
		sum += math.Pow(float64(s), 2)
	}

	// Square root
	return math.Sqrt(sum / float64(len(samples)))
}
