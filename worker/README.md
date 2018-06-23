# Worker

```go
// Init worker
var w = astiworker.NewWorker()

// Handle signals
w.HandleSignals()

// Serve
w.Serve("127.0.0.1:4000", myHandler)

// Execute
h, _ := w.Exec("sleep", "10")
go func() {
	time.Sleep(3 * time.Second)
	h.Stop()
}

// Wait
w.Wait()
```