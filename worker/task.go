package astiworker

import (
	"context"
	"sync"
)

// Task represents a task
type Task struct {
	ctx     context.Context
	od, ow  sync.Once
	wg, pwg *sync.WaitGroup
}

func newTask(parentCtx context.Context, parentWg *sync.WaitGroup) (t *Task) {
	t = &Task{
		wg:  &sync.WaitGroup{},
		pwg: parentWg,
	}
	t.ctx, _ = context.WithCancel(parentCtx)
	t.pwg.Add(1)
	return
}

// NewSubTask creates a new sub task
func (t *Task) NewSubTask() *Task {
	return newTask(t.ctx, t.wg)
}

// Ctx returns the task context
func (t *Task) Ctx() context.Context {
	return t.ctx
}

// Done indicates the task is done
func (t *Task) Done() {
	t.od.Do(func() {
		t.pwg.Done()
	})
}

// Wait waits for the task to be finished
// It is important to wait after the context is done
func (t *Task) Wait() {
	t.ow.Do(func() {
		<-t.ctx.Done()
		t.wg.Wait()
	})
}
