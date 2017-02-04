package astisync

import (
	"sync"
	"time"

	"fmt"
	"runtime"

	"github.com/rs/xlog"
)

// Constants
const (
	loggerKeyMutexName = "mutex_name"
)

// RWMutex represents a RWMutex capable of logging its actions to ease deadlock debugging
type RWMutex struct {
	lastSuccessfulLockCaller string
	logger                   xlog.Logger
	mutex                    *sync.RWMutex
	name                     string
}

// NewRWMutex creates a new RWMutex
func NewRWMutex(l xlog.Logger, name string) *RWMutex {
	return &RWMutex{
		logger: l,
		mutex:  &sync.RWMutex{},
		name:   name,
	}
}

// Lock write locks the mutex
func (m *RWMutex) Lock() {
	var caller string
	if _, file, line, ok := runtime.Caller(1); ok {
		caller = fmt.Sprintf("%s:%d", file, line)
	}
	m.logger.Debugf("Requesting lock for %s at %s", m.name, caller, xlog.F{
		loggerKeyMutexName: m.name,
	})
	m.mutex.Lock()
	m.logger.Debugf("Lock acquired for %s at %s", m.name, caller, xlog.F{
		loggerKeyMutexName: m.name,
	})
	m.lastSuccessfulLockCaller = caller
}

// Unlock write unlocks the mutex
func (m *RWMutex) Unlock() {
	m.mutex.Unlock()
	m.logger.Debugf("Unlock executed for %s", m.name, xlog.F{
		loggerKeyMutexName: m.name,
	})
}

// RLock read locks the mutex
func (m *RWMutex) RLock() {
	var caller string
	if _, file, line, ok := runtime.Caller(1); ok {
		caller = fmt.Sprintf("%s:%d", file, line)
	}
	m.logger.Debugf("Requesting rlock for %s at %s", m.name, caller, xlog.F{
		loggerKeyMutexName: m.name,
	})
	m.mutex.RLock()
	m.logger.Debugf("RLock acquired for %s at %s", m.name, caller, xlog.F{
		loggerKeyMutexName: m.name,
	})
	m.lastSuccessfulLockCaller = caller
}

// RUnlock read unlocks the mutex
func (m *RWMutex) RUnlock() {
	m.mutex.RUnlock()
	m.logger.Debugf("RUnlock executed for %s", m.name, xlog.F{
		loggerKeyMutexName: m.name,
	})
}

// IsDeadlocked checks whether the mutex is deadlocked with a given timeout and returns the last caller
func (m *RWMutex) IsDeadlocked(timeout time.Duration) (o bool, c string) {
	o = true
	c = m.lastSuccessfulLockCaller
	var channelLockAcquired = make(chan bool)
	go func() {
		m.mutex.Lock()
		defer m.mutex.Unlock()
		close(channelLockAcquired)
	}()
	for {
		select {
		case <-channelLockAcquired:
			o = false
			return
		case <-time.After(timeout):
			return
		}
	}
	return
}
