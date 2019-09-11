package astiaudio

import (
	"math"
	"sync"
	"time"
)

// SilenceDetector represents a silence detector
type SilenceDetector struct {
	analyses              []analysis
	buf                   []int32
	m                     *sync.Mutex // Locks buf
	minAnalysesPerSilence int
	o                     SilenceDetectorOptions
	samplesPerAnalysis    int
}

type analysis struct {
	level   float64
	samples []int32
}

// SilenceDetectorOptions represents a silence detector options
type SilenceDetectorOptions struct {
	MaxSilenceAudioLevel float64       `toml:"max_silence_audio_level"`
	MinSilenceDuration   time.Duration `toml:"min_silence_duration"`
	SampleRate           float64       `toml:"sample_rate"`
	StepDuration         time.Duration `toml:"step_duration"`
}

// NewSilenceDetector creates a new silence detector
func NewSilenceDetector(o SilenceDetectorOptions) (d *SilenceDetector) {
	// Create
	d = &SilenceDetector{
		m: &sync.Mutex{},
		o: o,
	}

	// Reset
	d.Reset()

	// Default option values
	if d.o.MinSilenceDuration == 0 {
		d.o.MinSilenceDuration = time.Second
	}
	if d.o.StepDuration == 0 {
		d.o.StepDuration = 30 * time.Millisecond
	}

	// Compute attributes depending on options
	d.samplesPerAnalysis = int(math.Floor(d.o.SampleRate * d.o.StepDuration.Seconds()))
	d.minAnalysesPerSilence = int(math.Floor(d.o.MinSilenceDuration.Seconds() / d.o.StepDuration.Seconds()))
	return
}

// Reset resets the silence detector
func (d *SilenceDetector) Reset() {
	// Lock
	d.m.Lock()
	defer d.m.Unlock()

	// Reset
	d.analyses = []analysis{}
	d.buf = []int32{}
}

// Add adds samples to the buffer and checks whether there are valid samples between silences
func (d *SilenceDetector) Add(samples []int32) (validSamples [][]int32) {
	// Lock
	d.m.Lock()
	defer d.m.Unlock()

	// Append samples to buffer
	d.buf = append(d.buf, samples...)

	// Analyze samples by step
	for len(d.buf) >= d.samplesPerAnalysis {
		// Append analysis
		d.analyses = append(d.analyses, analysis{
			level:   AudioLevel(d.buf[:d.samplesPerAnalysis]),
			samples: append([]int32(nil), d.buf[:d.samplesPerAnalysis]...),
		})

		// Remove samples from buffer
		d.buf = d.buf[d.samplesPerAnalysis:]
	}

	// Loop through analyses
	var leadingSilence, inBetween, trailingSilence int
	for i := 0; i < len(d.analyses); i++ {
		if d.analyses[i].level < d.o.MaxSilenceAudioLevel {
			// This is a silence

			// This is a leading silence
			if inBetween == 0 {
				leadingSilence++

				// The leading silence is valid
				// We can trim its useless part
				if leadingSilence > d.minAnalysesPerSilence {
					d.analyses = d.analyses[leadingSilence-d.minAnalysesPerSilence:]
					i -= leadingSilence - d.minAnalysesPerSilence
					leadingSilence = d.minAnalysesPerSilence
				}
				continue
			}

			// This is a trailing silence
			trailingSilence++

			// Trailing silence is invalid
			if trailingSilence < d.minAnalysesPerSilence {
				continue
			}

			// Trailing silence is valid
			// Loop through analyses
			var ss []int32
			for _, a := range d.analyses[:i+1] {
				ss = append(ss, a.samples...)
			}

			// Append valid samples
			validSamples = append(validSamples, ss)

			// Remove leading silence and non silence
			d.analyses = d.analyses[leadingSilence+inBetween:]
			i -= leadingSilence + inBetween

			// Reset counts
			leadingSilence, inBetween, trailingSilence = trailingSilence, 0, 0
		} else {
			// This is not a silence

			// This is a leading non silence
			// We need to remove it
			if i == 0 {
				d.analyses = d.analyses[1:]
				i = -1
				continue
			}

			// This is the first in-between
			if inBetween == 0 {
				// The leading silence is invalid
				// We need to remove it as well as this first non silence
				if leadingSilence < d.minAnalysesPerSilence {
					d.analyses = d.analyses[i+1:]
					i = -1
					continue
				}
			}

			// This non-silence was preceded by a silence not big enough to be a valid trailing silence
			// We incorporate it in the in-between
			if trailingSilence > 0 {
				inBetween += trailingSilence
				trailingSilence = 0
			}

			// This is an in-between
			inBetween++
			continue
		}
	}
	return
}
