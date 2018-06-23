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

// HandleSignals handles signals
func (w *Worker) HandleSignals() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch)
	go func() {
		for s := range ch {
			astilog.Infof("astiworker: received signal %s", s)
			if s == syscall.SIGABRT || s == syscall.SIGKILL || s == syscall.SIGINT || s == syscall.SIGQUIT || s == syscall.SIGTERM {
				w.Stop()
			}
		}
	}()
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
