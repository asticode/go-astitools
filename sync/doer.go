package astisync

import (
	"context"
	"sync"
)

// Doer is an object capable of doing stuff without blocking
type Doer struct {
	cancel context.CancelFunc
	cond   *sync.Cond
	ctx    context.Context
	fn     DoerFunc
	mc     *sync.Mutex // Locks cond
	mq     *sync.Mutex // Locks queue
	queue  []interface{}
}

// DoerFunc represents the doer func
type DoerFunc func(job interface{})

// NewDoer creates a new doer
func NewDoer(fn DoerFunc) (d *Doer) {
	// Create doer
	d = &Doer{
		fn: fn,
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
func (d *Doer) Close() error {
	d.cancel()
	return nil
}

// Do adds a job to the queue
func (d *Doer) Do(job interface{}) {
	// Add job
	d.mq.Lock()
	d.queue = append(d.queue, job)
	d.mq.Unlock()

	// Broadcast
	d.cond.L.Lock()
	d.cond.Broadcast()
	d.cond.L.Unlock()
}

// do loops through jobs and executes them if any, or wait for a new one otherwise
func (d *Doer) do() {
	for {
		// Check context
		if d.ctx.Err() != nil {
			return
		}

		// Lock cond here in case a job is added between retrieving l and doing the if on it
		d.cond.L.Lock()

		// Get number of jobs in queue
		d.mq.Lock()
		l := len(d.queue)
		d.mq.Unlock()

		// No queued job
		if l == 0 {
			d.cond.Wait()
			d.cond.L.Unlock()
			continue
		}
		d.cond.L.Unlock()

		// Get first job
		d.mq.Lock()
		job := d.queue[0]
		d.mq.Unlock()

		// Execute callback
		d.fn(job)

		// Remove first job
		d.mq.Lock()
		d.queue = d.queue[1:]
		d.mq.Unlock()
	}
}
