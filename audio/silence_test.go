package astiaudio

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSilenceDetector(t *testing.T) {
	// Create silence detector
	sd := NewSilenceDetector(SilenceDetectorOptions{
		MaxSilenceAudioLevel: 2,
		MinSilenceDuration:   400 * time.Millisecond, // 2 samples
		SampleRate:           5,
		StepDuration:         200 * time.Millisecond, // 1 sample
	})

	// Leading non silences + invalid leading silence + trailing silence is leftover
	vs := sd.Add([]int32{3, 1, 3, 1})
	assert.Equal(t, [][]int32(nil), vs)
	assert.Len(t, sd.analyses, 1)

	// Valid leading silence but trailing silence is insufficient for now
	vs = sd.Add([]int32{1, 3, 3, 1})
	assert.Equal(t, [][]int32(nil), vs)
	assert.Len(t, sd.analyses, 5)

	// Valid samples
	vs = sd.Add([]int32{1})
	assert.Equal(t, [][]int32{{1, 1, 3, 3, 1, 1}}, vs)
	assert.Len(t, sd.analyses, 2)

	// Multiple valid samples + truncate leading and trailing silences
	vs = sd.Add([]int32{1, 1, 1, 1, 3, 3, 1, 1, 1, 1, 3, 3, 1, 1, 1, 1})
	assert.Equal(t, [][]int32{{1, 1, 3, 3, 1, 1}, {1, 1, 3, 3, 1, 1}}, vs)
	assert.Len(t, sd.analyses, 2)

	// Invalid in-between silences that should be kept
	vs = sd.Add([]int32{1, 1, 1, 3, 3, 1, 3, 3, 1, 3, 3, 1, 1, 1})
	assert.Equal(t, [][]int32{{1, 1, 3, 3, 1, 3, 3, 1, 3, 3, 1, 1}}, vs)
	assert.Len(t, sd.analyses, 2)
}
