package astistat

import (
	"context"
	"sync"
	"time"
)

// Stater is an object that can compute and handle stats
type Stater struct {
	cancel context.CancelFunc
	ctx    context.Context
	fn     StatsHandleFunc
	oStart *sync.Once
	oStop  *sync.Once
	period time.Duration
	ss     []stat
}

// Stat represents a stat
type Stat struct {
	StatMetadata
	Value interface{}
}

// StatsHandleFunc is a method that can handle stats
type StatsHandleFunc func(stats []Stat)

// StatMetadata represents a stat metadata
type StatMetadata struct {
	Description string
	Label       string
	Unit        string
}

// StatValueFunc is a method that can compute a stat value
type StatValueFunc func(delta time.Duration) interface{}

// StatResetFunc is a method that can reset a stat
type StatResetFunc func()

type stat struct {
	fnReset StatResetFunc
	fnValue StatValueFunc
	m       StatMetadata
}

// NewStater creates a new stater
func NewStater(period time.Duration, fn StatsHandleFunc) *Stater {
	return &Stater{
		fn:     fn,
		oStart: &sync.Once{},
		oStop:  &sync.Once{},
		period: period,
	}
}

// Start starts the stater
func (s *Stater) Start(ctx context.Context) {
	// Make sure the stater can only be started once
	s.oStart.Do(func() {
		// Check context
		if ctx.Err() != nil {
			return
		}

		// Reset context
		s.ctx, s.cancel = context.WithCancel(ctx)

		// Reset once
		s.oStop = &sync.Once{}

		// Execute the rest in a go routine
		go func() {
			// Create ticker
			t := time.NewTicker(s.period)
			defer t.Stop()

			// Loop
			lastStatAt := time.Now()
			for {
				select {
				case <-t.C:
					// Get delta
					now := time.Now()
					delta := now.Sub(lastStatAt)
					lastStatAt = now

					// Loop through stats
					var stats []Stat
					for _, v := range s.ss {
						stats = append(stats, Stat{
							StatMetadata: v.m,
							Value:        v.fnValue(delta),
						})
					}

					// Handle stats
					go s.fn(stats)
				case <-s.ctx.Done():
					// Loop through stats
					for _, v := range s.ss {
						if v.fnReset != nil {
							v.fnReset()
						}
					}
					return
				}
			}
		}()
	})
}

// AddStat adds a stat
func (s *Stater) AddStat(m StatMetadata, fnValue StatValueFunc, fnReset StatResetFunc) {
	s.ss = append(s.ss, stat{
		fnReset: fnReset,
		fnValue: fnValue,
		m:       m,
	})
}

// Stop stops the stater
func (s *Stater) Stop() {
	// Make sure the stater can only be stopped once
	s.oStop.Do(func() {
		// Cancel context
		if s.cancel != nil {
			s.cancel()
		}

		// Reset once
		s.oStart = &sync.Once{}
	})
}

// StatsMetadata returns the stats metadata
func (s *Stater) StatsMetadata() (ms []StatMetadata) {
	ms = []StatMetadata{}
	for _, v := range s.ss {
		ms = append(ms, v.m)
	}
	return
}
