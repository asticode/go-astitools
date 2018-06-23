package astiworker

import (
	"context"
	"os/exec"
	"strings"
	"sync"

	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// Statuses
const (
	ExecHandlerStatusCrashed = "crashed"
	ExecHandlerStatusRunning = "running"
	ExecHandlerStatusStopped = "stopped"
)

// ExecHandler represents an object capable of handling the execution of a cmd
type ExecHandler interface {
	Status() string
	Stop()
}

type defaultExecHandler struct {
	cancel  context.CancelFunc
	ctx     context.Context
	err     error
	o       sync.Once
	stopped bool
}

func (h *defaultExecHandler) Status() string {
	if h.ctx.Err() != nil {
		if h.stopped || h.err == nil {
			return ExecHandlerStatusStopped
		}
		return ExecHandlerStatusCrashed
	}
	return ExecHandlerStatusRunning
}

func (h *defaultExecHandler) Stop() {
	h.o.Do(func() {
		h.cancel()
		h.stopped = true
	})
}

// Exec executes a cmd
// The process will be stopped when the worker stops
func (w *Worker) Exec(name string, args ...string) (ExecHandler, error) {
	// Create handler
	h := &defaultExecHandler{}
	h.ctx, h.cancel = context.WithCancel(context.Background())

	// Start
	cmd := exec.CommandContext(h.ctx, name, args...)
	n := strings.Join(append([]string{name}, args...), " ")
	astilog.Infof("astiworker: starting %s", n)
	if err := cmd.Start(); err != nil {
		err = errors.Wrapf(err, "astiworker: executing %s", n)
		return nil, err
	}

	// Make sure to increment the waiting group
	w.wg.Add(1)

	// Wait
	go func() {
		h.err = cmd.Wait()
		h.cancel()
		astilog.Infof("astiworker: status is now %s for %s", h.Status(), n)
		w.wg.Done()
	}()
	return h, nil
}
