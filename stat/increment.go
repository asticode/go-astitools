package astistat

import (
	"sync"
	"time"
)

// IncrementStat is an object capable of computing an increment stat properly
type IncrementStat struct {
	c         int64
	isStarted bool
	m         *sync.Mutex
}

// NewIncrementStat creates a new increment stat
func NewIncrementStat() *IncrementStat {
	return &IncrementStat{m: &sync.Mutex{}}
}

// Add increments the stat
func (s *IncrementStat) Add(delta int64) {
	s.m.Lock()
	defer s.m.Unlock()
	if !s.isStarted {
		return
	}
	s.c += delta
}

// Start implements the StatHandler interface
func (s *IncrementStat) Start() {
	s.m.Lock()
	defer s.m.Unlock()
	s.c = 0
	s.isStarted = true
}

// Stop implements the StatHandler interface
func (s *IncrementStat) Stop() {
	s.m.Lock()
	defer s.m.Unlock()
	s.isStarted = true
}

// Value implements the StatHandler interface
func (s *IncrementStat) Value(delta time.Duration) interface{} {
	s.m.Lock()
	defer s.m.Unlock()
	c := s.c
	s.c = 0
	return c
}
