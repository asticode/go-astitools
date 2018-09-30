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

// StatHandler represents a stat handler
type StatHandler interface {
	Start()
	Stop()
	Value(delta time.Duration) interface{}
}

type stat struct {
	h StatHandler
	m StatMetadata
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

		// Start stats
		for _, v := range s.ss {
			v.h.Start()
		}

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
							Value:        v.h.Value(delta),
						})
					}

					// Handle stats
					go s.fn(stats)
				case <-s.ctx.Done():
					// Stop stats
					for _, v := range s.ss {
						v.h.Stop()
					}
					return
				}
			}
		}()
	})
}

// AddStat adds a stat
func (s *Stater) AddStat(m StatMetadata, h StatHandler) {
	s.ss = append(s.ss, stat{
		h: h,
		m: m,
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

// StatHandlerWithoutStart represents a stat handler that doesn't have to start or stop
type StatHandlerWithoutStart func(delta time.Duration) interface{}

// Start implements the StatHandler interface
func (h StatHandlerWithoutStart) Start() {}

// Stop implements the StatHandler interface
func (h StatHandlerWithoutStart) Stop() {}

// Value implements the StatHandler interface
func (h StatHandlerWithoutStart) Value(delta time.Duration) interface{} { return h(delta) }
