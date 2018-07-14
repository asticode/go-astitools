package astiworker

import (
	"math"

	"github.com/asticode/go-astiamqp"
	"github.com/pkg/errors"
)

// ConfigurationConsumer represents a consumer configuration
type ConfigurationConsumer struct {
	AMQP        astiamqp.ConfigurationConsumer
	WorkerCount int
}

// Consume consumes AMQP events
func (w *Worker) Consume(a *astiamqp.AMQP, cs ...ConfigurationConsumer) (err error) {
	// Create task
	t := w.NewTask()

	// Loop through configurations
	for idxConf, c := range cs {
		// Loop through workers
		for idxWorker := 0; idxWorker < int(math.Max(1, float64(c.WorkerCount))); idxWorker++ {
			if err = a.AddConsumer(c.AMQP); err != nil {
				err = errors.Wrapf(err, "main: adding consumer #%d for conf #%d %+v failed", idxWorker+1, idxConf+1, c)
				return
			}
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
