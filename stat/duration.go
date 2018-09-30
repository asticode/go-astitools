package astistat

import (
	"sync"
	"time"
)

// DurationRatioStat is an object capable of computing a duration ratio stat properly
type DurationRatioStat struct {
	d         time.Duration
	isStarted bool
	m         *sync.Mutex
	startedAt map[interface{}]time.Time
}

// NewDurationRatioStat creates a new duration ratio stat
func NewDurationRatioStat() *DurationRatioStat {
	return &DurationRatioStat{
		startedAt: make(map[interface{}]time.Time),
		m:         &sync.Mutex{},
	}
}

// Add starts recording a new duration
func (s *DurationRatioStat) Add(k interface{}) {
	s.m.Lock()
	defer s.m.Unlock()
	if !s.isStarted {
		return
	}
	s.startedAt[k] = time.Now()
}

// Done indicates the duration is now done
func (s *DurationRatioStat) Done(k interface{}) {
	s.m.Lock()
	defer s.m.Unlock()
	if !s.isStarted {
		return
	}
	s.d += time.Now().Sub(s.startedAt[k])
	delete(s.startedAt, k)
}

// Value implements the StatHandler interface
func (s *DurationRatioStat) Value(delta time.Duration) (o interface{}) {
	// Lock
	s.m.Lock()
	defer s.m.Unlock()

	// Get current values
	n := time.Now()
	d := s.d

	// Loop through waits still not finished
	for k, v := range s.startedAt {
		d += n.Sub(v)
		s.startedAt[k] = n
	}

	// Compute stat
	o = float64(d) / float64(delta) * 100
	s.d = 0
	return
}

// Start implements the StatHandler interface
func (s *DurationRatioStat) Start() {
	s.m.Lock()
	defer s.m.Unlock()
	s.d = 0
	s.isStarted = true
	s.startedAt = make(map[interface{}]time.Time)
}

// Stop implements the StatHandler interface
func (s *DurationRatioStat) Stop() {
	s.m.Lock()
	defer s.m.Unlock()
	s.isStarted = false
}
