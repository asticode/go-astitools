package astistat

import (
	"sync"
	"time"
)

// WaitStat is an object capable of computing a wait stat properly
type WaitStat struct {
	d         time.Duration
	isStarted bool
	m         *sync.Mutex
	startedAt map[interface{}]time.Time
}

// NewWaitStat creates a new wait stat
func NewWaitStat() *WaitStat {
	return &WaitStat{
		startedAt: make(map[interface{}]time.Time),
		m:         &sync.Mutex{},
	}
}

// Add adds a new wait
func (s *WaitStat) Add(k interface{}) {
	s.m.Lock()
	defer s.m.Unlock()
	if !s.isStarted {
		return
	}
	s.startedAt[k] = time.Now()
}

// WaitFinished indicates a wait has finished
func (s *WaitStat) Done(k interface{}) {
	s.m.Lock()
	defer s.m.Unlock()
	if !s.isStarted {
		return
	}
	s.d += time.Now().Sub(s.startedAt[k])
	delete(s.startedAt, k)
}

// Value implements the StatHandler interface
func (s *WaitStat) Value(delta time.Duration) (o interface{}) {
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
func (s *WaitStat) Start() {
	s.m.Lock()
	defer s.m.Unlock()
	s.d = 0
	s.isStarted = true
	s.startedAt = make(map[interface{}]time.Time)
}

// Stop implements the StatHandler interface
func (s *WaitStat) Stop() {
	s.m.Lock()
	defer s.m.Unlock()
	s.isStarted = false
}
