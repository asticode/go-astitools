package astisync

import (
	"context"
	"sync"
)

// CtxQueue is a queue that can handle a context without dropping any message in between
type CtxQueue struct {
	c         chan ctxQueueMessage
	ctxIsDone bool
	m         *sync.Mutex
	o         *sync.Once
}

type ctxQueueMessage struct {
	ctxIsDone bool
	p         interface{}
}

// NewCtxQueue creates a new ctx queue
func NewCtxQueue() *CtxQueue {
	return &CtxQueue{
		c: make(chan ctxQueueMessage),
		m: &sync.Mutex{},
		o: &sync.Once{},
	}
}

// Start starts the queue
func (q *CtxQueue) Start(ctx context.Context, fn func(p interface{})) {
	// Make sure the queue can only be started once
	q.o.Do(func() {
		// Handle ctx
		go q.handleCtx(ctx)

		// Loop
		for {
			select {
			case m := <-q.c:
				// Check context
				if m.ctxIsDone {
					return
				}

				// Handle payload
				fn(m.p)
			}
		}
	})
}

func (q *CtxQueue) handleCtx(ctx context.Context) {
	<-ctx.Done()
	q.c <- ctxQueueMessage{ctxIsDone: true}
}

// Send sends a message in the queue
func (q *CtxQueue) Send(p interface{}) {
	q.c <- ctxQueueMessage{p: p}
}
