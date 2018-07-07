package astiworker

import (
	"math"

	"github.com/asticode/go-astiamqp"
	"github.com/pkg/errors"
)

// Consume consumes AMQP events
func (w *Worker) Consume(a *astiamqp.AMQP, c astiamqp.ConfigurationConsumer, workerCount int) (err error) {
	// Create task
	t := w.NewTask()

	// Add consumers
	for idx := 0; idx < int(math.Max(1, float64(workerCount))); idx++ {
		if err = a.AddConsumer(c); err != nil {
			err = errors.Wrapf(err, "main: adding consumer #%d with conf %+v failed", idx+1, c)
			return
		}
	}

	// Execute the rest in a go routine
	go func() {
		// Wait for context to be done
		<-w.Context().Done()

		// Stop amqp
		a.Stop()

		// Task is done
		t.Done()
	}()
	return
}
