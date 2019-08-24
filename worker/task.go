package astiworker

import (
	"sync"
)

// Task represents a task
type Task struct {
	od, ow  sync.Once
	wg, pwg *sync.WaitGroup
}

func newTask(parentWg *sync.WaitGroup) (t *Task) {
	t = &Task{
		wg:  &sync.WaitGroup{},
		pwg: parentWg,
	}
	t.pwg.Add(1)
	return
}

// TaskFunc represents a function that can create a new task
type TaskFunc func() *Task

// NewSubTask creates a new sub task
func (t *Task) NewSubTask() *Task {
	return newTask(t.wg)
}

// Done indicates the task is done
func (t *Task) Done() {
	t.od.Do(func() {
		t.pwg.Done()
	})
}

// Wait waits for the task to be finished
func (t *Task) Wait() {
	t.ow.Do(func() {
		t.wg.Wait()
	})
}
