package astitime_test

import (
	"testing"
	stltime "time"

	"sync"

	"github.com/asticode/go-toolkit/time"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestSleep(t *testing.T) {
	var ctx, cancel = context.WithCancel(context.Background())
	var err error
	var wg = &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = time.Sleep(ctx, stltime.Minute)
	}()
	cancel()
	wg.Wait()
	assert.EqualError(t, err, "context canceled")
}
