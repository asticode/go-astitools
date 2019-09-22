package astievent

import (
	"context"
	"sync"

	astisync "github.com/asticode/go-astitools/sync"
)

// Dispatcher represents an object that can dispatch simple events (name + payload)
type Dispatcher struct {
	c  *astisync.Chan
	hs map[string][]EventHandler
	mh *sync.Mutex
}

// EventHandler represents a function that can handler an event's payload
type EventHandler func(payload interface{})

// NewDispatcher creates a new dispatcher
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		c:  astisync.NewChan(astisync.ChanOptions{}),
		hs: make(map[string][]EventHandler),
		mh: &sync.Mutex{},
	}
}

// On adds an event handler for a specific name
func (d *Dispatcher) On(name string, h EventHandler) {
	// Lock
	d.mh.Lock()
	defer d.mh.Unlock()

	// Add handler
	d.hs[name] = append(d.hs[name], h)
}

// Dispatch dispatches a payload for a specific name
func (d *Dispatcher) Dispatch(name string, payload interface{}) {
	// Lock
	d.mh.Lock()
	defer d.mh.Unlock()

	// No handlers
	hs, ok := d.hs[name]
	if !ok {
		return
	}

	// Loop through handlers
	for _, h := range hs {
		// We need to store the handler
		sh := h

		// Add to chan
		d.c.Add(func() {
			sh(payload)
		})
	}
}

// Start starts the dispatcher. It is blocking
func (d *Dispatcher) Start(ctx context.Context) {
	d.c.Start(ctx)
}

// Stop stops the dispatcher
func (d *Dispatcher) Stop() {
	d.c.Stop()
}

// Reset resets the dispatcher
func (d *Dispatcher) Reset() {
	d.c.Reset()
}
