package astiworker

import (
	"os"
	"os/signal"
	"syscall"

	"sync"

	"github.com/asticode/go-astilog"
)

// Worker represents an object capable of blocking, handling signals and stopping
type Worker struct {
	channelQuit chan bool
	os          sync.Once
	ow          sync.Once
}

// NewWorker builds a new worker
func NewWorker() *Worker {
	astilog.Info("astiworker: starting worker...")
	return &Worker{channelQuit: make(chan bool)}
}

// HandleSignals handles signals
func (w Worker) HandleSignals() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch)
	go func() {
		for s := range ch {
			astilog.Infof("astiworker: received signal %s", s)
			if s == syscall.SIGABRT || s == syscall.SIGKILL || s == syscall.SIGINT || s == syscall.SIGQUIT || s == syscall.SIGTERM {
				w.Stop()
			}
			return
		}
	}()
}

// Stop stops the Worker
func (w *Worker) Stop() {
	w.os.Do(func() {
		astilog.Info("astiworker: stopping worker...")
		close(w.channelQuit)
	})
}

// Wait is a blocking pattern
func (w *Worker) Wait() {
	w.ow.Do(func() {
		astilog.Info("astiworker: worker is now waiting...")
		<-w.channelQuit
	})
}
