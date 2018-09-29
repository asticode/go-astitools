package astistat

import (
	"sync"
	"time"
)

// WaitStat is an object capable of computing a wait stat properly
type WaitStat struct {
	startedAt map[interface{}]time.Time
	d         time.Duration
	m         *sync.Mutex
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
	s.startedAt[k] = time.Now()
}

// WaitFinished indicates a wait has finished
func (s *WaitStat) Done(k interface{}) {
	s.m.Lock()
	defer s.m.Unlock()
	s.d += time.Now().Sub(s.startedAt[k])
	delete(s.startedAt, k)
}

// StatValueFunc is the wait stat value func
func (s *WaitStat) StatValueFunc(delta time.Duration) (o interface{}) {
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

// Reset resets the stat
func (s *WaitStat) Reset() {
	s.d = 0
}
