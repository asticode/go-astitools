package astisync

import (
	"context"
	"sync"
	"sync/atomic"
)

// Regulator is an object that can keep track of processes which, in turn, keep track of subprocesses.
// It ensures 2 things:
//   - if a limit [n] is provided, no more than [n] processes are running simultaneously
//   - when a process runs, it only gives control back to the main loop when one of its subprocess has finished
type Regulator struct {
	c              *sync.Cond
	limit          int
	parentCtx      context.Context
	processesCount int
	subprocessesWg *sync.WaitGroup
}

// NewRegulator creates a new regulator
func NewRegulator(parentCtx context.Context, limit int) *Regulator {
	return &Regulator{
		c:              sync.NewCond(&sync.Mutex{}),
		limit:          limit,
		parentCtx:      parentCtx,
		subprocessesWg: &sync.WaitGroup{},
	}
}

// NewProcess creates a new regulator process
func (r *Regulator) NewProcess() *RegulatorProcess {
	// Check the number of process already running
	// If this number is bigger than the limit we wait for one of the process to finish
	r.c.L.Lock()
	if r.limit > 0 && r.processesCount >= r.limit {
		r.c.Wait()
	}

	// Create process
	r.processesCount++
	r.c.L.Unlock()
	return newRegulatorProcess(r.parentCtx, r.processIsDone, r.subprocessesWg)
}

// Wait waits for all subprocesses to be finished
func (r *Regulator) Wait() {
	r.subprocessesWg.Wait()
}

// This method is called by a process when it is done
// It will decrement the number of running processes and broadcast the fact that a process slot is now available
func (r *Regulator) processIsDone() {
	r.c.L.Lock()
	r.processesCount--
	r.c.Broadcast()
	r.c.L.Unlock()
}

// RegulatorProcess is a regulator process
type RegulatorProcess struct {
	c                 *sync.Cond
	cancel            context.CancelFunc
	ctx               context.Context
	doneFuncs         []func()
	subprocessesWg    *sync.WaitGroup
	subprocessesCount int64
}

func newRegulatorProcess(ctx context.Context, doneFunc func(), subprocessesWg *sync.WaitGroup) (p *RegulatorProcess) {
	p = &RegulatorProcess{
		c:              sync.NewCond(&sync.Mutex{}),
		doneFuncs:      []func(){doneFunc},
		subprocessesWg: subprocessesWg,
	}
	p.c.L.Lock()
	p.ctx, p.cancel = context.WithCancel(ctx)
	return
}

// AddSubprocesses adds subprocesses to the process
func (p *RegulatorProcess) AddSubprocesses(delta int, doneFunc func()) {
	// Make sure we keep track of all subprocesses
	p.subprocessesWg.Add(delta)

	// Increment the number of running subprocesses
	atomic.AddInt64(&p.subprocessesCount, int64(delta))

	// Add done func
	if doneFunc != nil {
		p.doneFuncs = append([]func(){doneFunc}, p.doneFuncs...)
	}
}

// Wait waits either for one of the children to be finished or for the ctx to be cancelled
func (p *RegulatorProcess) Wait() {
	// Listen to context
	go func() {
		<-p.ctx.Done()
		p.broadcast()
	}()

	// Wait for a broadcast
	p.c.Wait()
	p.c.L.Unlock()
}

// SubprocessIsDone indicates that one of the subprocess is done
func (p *RegulatorProcess) SubprocessIsDone() {
	// Broadcast
	p.broadcast()

	// Make sure to keep track of all subprocesses
	p.subprocessesWg.Done()

	// Decrement the number of running subprocesses
	if c := atomic.AddInt64(&p.subprocessesCount, -1); c == 0 {
		// Execute the done funcs
		for _, f := range p.doneFuncs {
			f()
		}

		// Cancel the context so the goroutine is closed
		p.cancel()
	}
}

func (p *RegulatorProcess) broadcast() {
	p.c.L.Lock()
	p.c.Broadcast()
	p.c.L.Unlock()
}
