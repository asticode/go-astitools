package astisync

import (
	"context"
	"sync"
)

// Do is an object capable of doing stuff in FIFO order without blocking
type Do struct {
	cancel context.CancelFunc
	cond   *sync.Cond
	ctx    context.Context
	mc     *sync.Mutex // Locks cond
	mq     *sync.Mutex // Locks queue
	queue  []func()
}

// NewDo creates a new Do
func NewDo() (d *Do) {
	// Create do
	d = &Do{
		mc: &sync.Mutex{},
		mq: &sync.Mutex{},
	}

	// Create cond
	d.cond = sync.NewCond(d.mc)

	// Create context
	d.ctx, d.cancel = context.WithCancel(context.Background())

	// Do
	go d.do()
	return
}

// Close implements the io.Closer interface
func (d *Do) Close() error {
	d.cancel()
	return nil
}

// Do execute a new func
func (d *Do) Do(fn func()) {
	// Add job
	d.mq.Lock()
	d.queue = append(d.queue, fn)
	d.mq.Unlock()

	// Broadcast
	d.cond.L.Lock()
	d.cond.Broadcast()
	d.cond.L.Unlock()
}

// do loops through funcs in queue and executes them if any, or wait for a new one otherwise
func (d *Do) do() {
	for {
		// Check context
		if d.ctx.Err() != nil {
			return
		}

		// Lock cond here in case a func is added between retrieving l and doing the if on it
		d.cond.L.Lock()

		// Get number of funcs in queue
		d.mq.Lock()
		l := len(d.queue)
		d.mq.Unlock()

		// No queued funcs
		if l == 0 {
			d.cond.Wait()
			d.cond.L.Unlock()
			continue
		}
		d.cond.L.Unlock()

		// Get first func
		d.mq.Lock()
		fn := d.queue[0]
		d.mq.Unlock()

		// Execute func
		fn()

		// Remove first func
		d.mq.Lock()
		d.queue = d.queue[1:]
		d.mq.Unlock()
	}
}
