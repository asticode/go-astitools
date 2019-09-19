package astipcm

import (
	"math"
)

func maxSample(bitDepth int) int {
	return int(math.Pow(2, float64(bitDepth))/2.0) - 1
}

func SampleToAmplitude(s int, bitDepth int, signed bool) (a float64) {
	// Get max
	max := maxSample(bitDepth)

	// Compute amplitude
	if signed {
		// Sample values are between -max <= x <= max
		// We need them to be 0 <= x <= 1 so we first make them be -0.5 <= x <= 0.5 and we add 0.5
		a = (float64(s) / (float64(max) * 2.0)) + 0.5
	} else {
		a = float64(s) / float64(max)
	}
	return
}

func AmplitudeToSample(a float64, bitDepth int, signed bool) (s int) {
	// Get max
	max := maxSample(bitDepth)

	// Compute sample
	if signed {
		s = int((a - 0.5) * (float64(max) * 2.0))
	} else {
		s = int(a * float64(max))
	}
	return
}

func AmplitudeToDB(a float64) float64 {
	return 20 * math.Log10(a)
}

func DBToAmplitude(db float64) float64 {
	return math.Pow(10.0, db*0.05)
}

// https://stackoverflow.com/questions/2445756/how-can-i-calculate-audio-db-level
func SampleToDB(s int, bitDepth int, signed bool) float64 {
	return AmplitudeToDB(SampleToAmplitude(s, bitDepth, signed))
}

func DBToSample(db float64, bitDepth int, signed bool) int {
	return AmplitudeToSample(DBToAmplitude(db), bitDepth, signed)
}

func Normalize(samples []int, bitDepth int) (o []int) {
	// Get max sample
	var m int
	for _, s := range samples {
		if v := int(math.Abs(float64(s))); v > m {
			m = v
		}
	}

	// Get max for bit depth
	max := maxSample(bitDepth)

	// Loop through samples
	for _, s := range samples {
		o = append(o, s*max/m)
	}
	return
}
