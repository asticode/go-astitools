package astihttp

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astitools/io"
	"github.com/pkg/errors"
)

// Downloader represents a downloader
type Downloader struct {
	busyWorkers     int
	c               *http.Client
	cond            *sync.Cond
	mc              *sync.Mutex // Locks cond
	mw              *sync.Mutex // Locks busyWorkers
	numberOfWorkers int
}

// DownloaderFunc represents a downloader func
type DownloaderFunc func(ctx context.Context, idx int, src string, buf *bytes.Buffer) error

// DownloaderOptions represents downloader options
type DownloaderOptions struct {
	Client          *http.Client
	NumberOfWorkers int
}

// NewDownloader creates a new downloader
func NewDownloader(o DownloaderOptions) (d *Downloader) {
	d = &Downloader{
		c:               o.Client,
		mc:              &sync.Mutex{},
		mw:              &sync.Mutex{},
		numberOfWorkers: o.NumberOfWorkers,
	}
	if d.c == nil {
		d.c = &http.Client{}
	}
	d.cond = sync.NewCond(d.mc)
	if d.numberOfWorkers == 0 {
		d.numberOfWorkers = 1
	}
	return
}

// Download downloads in parallel a set of src paths and executes a custom callback on each downloaded buffers
func (d *Downloader) Download(ctx context.Context, srcs []string, fn DownloaderFunc) (err error) {
	// Loop through src paths
	wg := &sync.WaitGroup{}
	wg.Add(len(srcs))
	var idx int
	for idx < len(srcs) {
		// Check context
		if ctx.Err() != nil {
			err = errors.Wrap(err, "astihttp: context error")
			return
		}

		// Lock cond here in case a worker finishes between checking the number of busy workers and the if statement
		d.cond.L.Lock()

		// Check if a worker is available
		var ok bool
		d.mw.Lock()
		if ok = d.numberOfWorkers > d.busyWorkers; ok {
			d.busyWorkers++
		}
		d.mw.Unlock()

		// No worker is available
		if !ok {
			d.cond.Wait()
			d.cond.L.Unlock()
			continue
		}
		d.cond.L.Unlock()

		// Download
		go func(idx int) {
			if errR := d.download(ctx, idx, srcs[idx], fn, wg); errR != nil {
				err = errR
			}
		}(idx)
		idx++
	}
	wg.Wait()
	return
}

func (d *Downloader) download(ctx context.Context, idx int, src string, fn DownloaderFunc, wg *sync.WaitGroup) (err error) {
	// Update wait group and worker status
	defer func() {
		// Update worker status
		d.mw.Lock()
		d.busyWorkers--
		d.mw.Unlock()

		// Broadcast
		d.cond.L.Lock()
		d.cond.Broadcast()
		d.cond.L.Unlock()

		// Update wait group
		wg.Done()
	}()

	// Download
	buf := &bytes.Buffer{}
	astilog.Debugf("astihttp: downloading %s", src)
	if err = DownloadInWriter(ctx, d.c, src, buf); err != nil {
		err = errors.Wrapf(err, "astihttp: downloading %s failed", src)
		return
	}

	// Custom callback
	if err = fn(ctx, idx, src, buf); err != nil {
		err = errors.Wrapf(err, "astihttp: custom callback on %s failed", src)
		return
	}
	return
}

// DownloadInDirectory downloads in parallel a set of src paths and saves them in a dst directory
func (d *Downloader) DownloadInDirectory(ctx context.Context, dst string, srcs ...string) error {
	return d.Download(ctx, srcs, func(ctx context.Context, idx int, src string, buf *bytes.Buffer) (err error) {
		// Make sure destination directory exists
		if err = os.MkdirAll(dst, 0700); err != nil {
			err = errors.Wrapf(err, "astihttp: mkdirall %s failed", dst)
			return
		}

		// Create destination file
		var f *os.File
		dst := filepath.Join(dst, filepath.Base(src))
		if f, err = os.Create(dst); err != nil {
			err = errors.Wrapf(err, "astihttp: creating %s failed", dst)
			return
		}
		defer f.Close()

		// Copy
		if _, err = astiio.Copy(ctx, buf, f); err != nil {
			err = errors.Wrapf(err, "astihttp: copying content to %s failed", dst)
			return
		}
		return
	})
}
