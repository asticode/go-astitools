package astiworker

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/asticode/go-astilog"
)

// Worker represents an object capable of blocking, handling signals and stopping
type Worker struct {
	cancel context.CancelFunc
	ctx    context.Context
	os, ow sync.Once
	wg     *sync.WaitGroup
}

// NewWorker builds a new worker
func NewWorker() (w *Worker) {
	astilog.Info("astiworker: starting worker...")
	w = &Worker{wg: &sync.WaitGroup{}}
	w.ctx, w.cancel = context.WithCancel(context.Background())
	w.wg.Add(1)
	return
}

// SignalHandler represents a func that can handle a signal
type SignalHandler func(s os.Signal)

func isTermSignal(s os.Signal) bool {
	return s == syscall.SIGABRT || s == syscall.SIGKILL || s == syscall.SIGINT || s == syscall.SIGQUIT || s == syscall.SIGTERM
}

// TermSignalHandler returns a SignalHandler that is executed only on a term signal
func TermSignalHandler(f func()) SignalHandler {
	return func(s os.Signal) {
		if isTermSignal(s) {
			f()
		}
	}
}

// HandleSignals handles signals
func (w *Worker) HandleSignals(hs ...SignalHandler) {
	// Add default handler
	hs = append([]SignalHandler{TermSignalHandler(w.Stop)}, hs...)

	// Notify
	ch := make(chan os.Signal, 1)
	signal.Notify(ch)

	// Execute in a task
	w.NewTask().Do(func() {
		for {
			select {
			case s := <-ch:
				// Log
				astilog.Debugf("astiworker: received signal %s", s)

				// Loop through handlers
				for _, h := range hs {
					h(s)
				}

				// Return
				if isTermSignal(s) {
					return
				}
			case <-w.Context().Done():
				return
			}
		}
	})
}

// Stop stops the Worker
func (w *Worker) Stop() {
	w.os.Do(func() {
		astilog.Info("astiworker: stopping worker...")
		w.cancel()
		w.wg.Done()
	})
}

// Wait is a blocking pattern
func (w *Worker) Wait() {
	w.ow.Do(func() {
		astilog.Info("astiworker: worker is now waiting...")
		w.wg.Wait()
	})
}

// NewTask creates a new task
func (w *Worker) NewTask() *Task {
	return newTask(w.wg)
}

// Context returns the worker's context
func (w *Worker) Context() context.Context {
	return w.ctx
}
