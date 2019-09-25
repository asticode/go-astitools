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
	StatusCrashed = "crashed"
	StatusRunning = "running"
	StatusStopped = "stopped"
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
			return StatusStopped
		}
		return StatusCrashed
	}
	return StatusRunning
}

func (h *defaultExecHandler) Stop() {
	h.o.Do(func() {
		h.cancel()
		h.stopped = true
	})
}

type ExecOptions struct {
	Args       []string
	CmdAdapter func(cmd *exec.Cmd) error
	Name       string
}

// Exec executes a cmd
// The process will be stopped when the worker stops
func (w *Worker) Exec(o ExecOptions) (ExecHandler, error) {
	// Create handler
	h := &defaultExecHandler{}
	h.ctx, h.cancel = context.WithCancel(context.Background())

	// Create command
	cmd := exec.CommandContext(h.ctx, o.Name, o.Args...)

	// Adapt command
	if o.CmdAdapter != nil {
		if err := o.CmdAdapter(cmd); err != nil {
			return nil, errors.Wrap(err, "astiworker: adapting cmd failed")
		}
	}

	// Start
	astilog.Infof("astiworker: starting %s", strings.Join(cmd.Args, " "))
	if err := cmd.Start(); err != nil {
		err = errors.Wrapf(err, "astiworker: executing %s", strings.Join(cmd.Args, " "))
		return nil, err
	}

	// Create task
	t := w.NewTask()

	// Wait
	go func() {
		h.err = cmd.Wait()
		h.cancel()
		astilog.Infof("astiworker: status is now %s for %s", h.Status(), strings.Join(cmd.Args, " "))
		t.Done()
	}()
	return h, nil
}
