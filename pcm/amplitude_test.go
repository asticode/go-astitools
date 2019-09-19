package astipcm

import (
	"fmt"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxSample(t *testing.T) {
	assert.Equal(t, 127, maxSample(8))
	assert.Equal(t, 32767, maxSample(16))
	assert.Equal(t, 8388607, maxSample(24))
	assert.Equal(t, 2147483647, maxSample(32))
}

func TestSampleToAmplitude(t *testing.T) {
	// Signed boundaries
	assert.Equal(t, 1.0, SampleToAmplitude(maxSample(16), 16, true))
	assert.Equal(t, 0.0, SampleToAmplitude(-maxSample(16), 16, true))

	// Signed value
	assert.Equal(t, 0.5, SampleToAmplitude(0, 16, true))

	// Unsigned boundaries
	assert.Equal(t, 1.0, SampleToAmplitude(maxSample(16), 16, false))
	assert.Equal(t, 0.0, SampleToAmplitude(0, 16, false))

	// Unsigned value
	assert.Equal(t, 0.499984740745262, SampleToAmplitude(maxSample(16)/2, 16, false))
}

func TestAmplitudeToSample(t *testing.T) {
	// Signed boundaries
	assert.Equal(t, maxSample(16), AmplitudeToSample(1.0, 16, true))
	assert.Equal(t, -maxSample(16), AmplitudeToSample(0.0, 16, true))

	// Signed value
	assert.Equal(t, 0, AmplitudeToSample(0.5, 16, true))

	// Unsigned boundaries
	assert.Equal(t, maxSample(16), AmplitudeToSample(1.0, 16, false))
	assert.Equal(t, 0, AmplitudeToSample(0.0, 16, false))

	// Unsigned value
	assert.Equal(t, maxSample(16)/2.0, AmplitudeToSample(0.499984740745262, 16, false))
}

func TestAmplitudeToDB(t *testing.T) {
	assert.Equal(t, math.Inf(-1), AmplitudeToDB(0.0))
	assert.Equal(t, -7.1309464702762515, AmplitudeToDB(0.44))
	assert.Equal(t, 0.0, AmplitudeToDB(1.0))
}

func TestDBToAmplitude(t *testing.T) {
	assert.Equal(t, 0.0, DBToAmplitude(math.Inf(-1)))
	assert.Equal(t, 0.44004794783598367, DBToAmplitude(-7.13))
	assert.Equal(t, 1.0, DBToAmplitude(0.0))
}

func TestSampleToDB(t *testing.T) {
	assert.Equal(t, -6.020599913279624, SampleToDB(0, 16, true))
	assert.Equal(t, -2.4988630927505144, SampleToDB(maxSample(16)/2, 16, true))
}

func TestDBToSample(t *testing.T) {
	assert.Equal(t, 0, DBToSample(-6.020599913279624, 16, true))
	assert.Equal(t, maxSample(16)/2, DBToSample(-2.4988630927505144, 16, true))
}

func TestNormalize(t *testing.T) {
	// Nothing to do
	i := []int{10000, maxSample(16), -10000}
	assert.Equal(t, i, Normalize(i, 16))

	fmt.Fprintln(os.Stdout, "")

	// Normalize
	i = []int{10000, 0, -10000}
	assert.Equal(t, []int{32767, 0, -32767}, Normalize(i, 16))
}
