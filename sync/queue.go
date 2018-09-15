package astisync

import (
	"context"
	"sync"
	"sync/atomic"
)

// CtxQueue is a queue that can handle a context without dropping any message in between
type CtxQueue struct {
	c         chan ctxQueueMessage
	ctxIsDone uint32
	o         *sync.Once
}

type ctxQueueMessage struct {
	c         *sync.Cond
	ctxIsDone bool
	p         interface{}
}

// NewCtxQueue creates a new ctx queue
func NewCtxQueue() *CtxQueue {
	return &CtxQueue{
		c: make(chan ctxQueueMessage),
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

				// Broadcast the fact that the process is done
				if m.c != nil {
					m.c.L.Lock()
					m.c.Broadcast()
					m.c.L.Unlock()
				}
			}
		}
	})
}

func (q *CtxQueue) handleCtx(ctx context.Context) {
	<-ctx.Done()
	atomic.StoreUint32(&q.ctxIsDone, 1)
	q.c <- ctxQueueMessage{ctxIsDone: true}
}

// Send sends a message in the queue
func (q *CtxQueue) Send(p interface{}) {
	// Context is done
	if d := atomic.LoadUint32(&q.ctxIsDone); d == 1 {
		return
	}

	// Send message
	q.c <- ctxQueueMessage{p: p}
}

// SendAndWait sends a message in the queue and waits for the end of its handling
func (q *CtxQueue) SendAndWait(p interface{}) {
	// Context is done
	if d := atomic.LoadUint32(&q.ctxIsDone); d == 1 {
		return
	}

	// Create cond
	c := sync.NewCond(&sync.Mutex{})
	c.L.Lock()

	// Send message
	q.c <- ctxQueueMessage{
		c: c,
		p: p,
	}

	// Wait for handling to be done
	c.Wait()
}
