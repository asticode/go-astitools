package astiworker

import (
	"os"
	"os/signal"
	"syscall"

	"context"

	"sync"

	"net/http"

	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
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

// Serve spawns a server
func (w *Worker) Serve(addr string, h http.Handler) {
	// Create server
	s := &http.Server{Addr: addr, Handler: h}

	// Make sure to increment the waiting group
	w.wg.Add(1)

	// Execute the rest in a goroutine
	astilog.Infof("astiworker: serving on %s", addr)
	go func() {
		// Serve
		var chanDone = make(chan error)
		go func() {
			if err := s.ListenAndServe(); err != nil {
				chanDone <- err
			}
		}()

		// Wait for context or chanDone to be done
		select {
		case <-w.ctx.Done():
			if w.ctx.Err() != context.Canceled {
				astilog.Error(errors.Wrap(w.ctx.Err(), "astiworker: context error"))
			}
		case err := <-chanDone:
			if err != nil {
				astilog.Error(errors.Wrap(err, "astiworker: serving failed"))
			}
		}

		// Shutdown
		astilog.Debugf("astiworker: shutting down server on %s", addr)
		if err := s.Shutdown(context.Background()); err != nil {
			astilog.Error(errors.Wrapf(err, "astiworker: shutting down server on %s failed", addr))
		}
		w.wg.Done()
	}()
}
