package astihttp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/asticode/go-astitools/io"
	"github.com/pkg/errors"
)

// Downloader represents a downloader
type Downloader struct {
	busyWorkers     int
	cond            *sync.Cond
	mc              *sync.Mutex // Locks cond
	mw              *sync.Mutex // Locks busyWorkers
	numberOfWorkers int
	s               *Sender
}

// DownloaderFunc represents a downloader func
// It's its responsibility to close the reader
type DownloaderFunc func(ctx context.Context, idx int, src string, r io.ReadCloser) error

// DownloaderOptions represents downloader options
type DownloaderOptions struct {
	NumberOfWorkers int
	Sender          SenderOptions
}

// NewDownloader creates a new downloader
func NewDownloader(o DownloaderOptions) (d *Downloader) {
	d = &Downloader{
		mc:              &sync.Mutex{},
		mw:              &sync.Mutex{},
		numberOfWorkers: o.NumberOfWorkers,
		s:               NewSender(o.Sender),
	}
	d.cond = sync.NewCond(d.mc)
	if d.numberOfWorkers == 0 {
		d.numberOfWorkers = 1
	}
	return
}

// Download downloads in parallel a set of src paths and executes a custom callback on each downloaded buffers
func (d *Downloader) Download(ctx context.Context, paths []string, fn DownloaderFunc) (err error) {
	// Loop through src paths
	m := &sync.Mutex{} // Locks err
	wg := &sync.WaitGroup{}
	wg.Add(len(paths))
	var idx int
	for idx < len(paths) {
		// Check context
		if ctx.Err() != nil {
			m.Lock()
			err = errors.Wrap(err, "astihttp: context error")
			m.Unlock()
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

		// Check error
		m.Lock()
		if err != nil {
			m.Unlock()
			return
		}
		m.Unlock()

		// Download
		go func(idx int) {
			if errR := d.download(ctx, idx, paths[idx], fn, wg); errR != nil {
				m.Lock()
				if err == nil {
					err = errR
				}
				m.Unlock()
			}
		}(idx)
		idx++
	}
	wg.Wait()
	return
}

func (d *Downloader) download(ctx context.Context, idx int, path string, fn DownloaderFunc, wg *sync.WaitGroup) (err error) {
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

	// Create request
	var r *http.Request
	if r, err = http.NewRequest(http.MethodGet, path, nil); err != nil {
		return errors.Wrapf(err, "astihttp: creating GET request to %s failed", path)
	}

	// Send request
	var resp *http.Response
	if resp, err = d.s.Send(r); err != nil {
		return errors.Wrapf(err, "astihttp: sending GET request to %s failed", path)
	}

	// Validate status code
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return fmt.Errorf("astihttp: sending GET request to %s returned %d status code", path, resp.StatusCode)
	}

	// Custom callback
	if err = fn(ctx, idx, path, resp.Body); err != nil {
		err = errors.Wrapf(err, "astihttp: custom callback on %s failed", path)
		return
	}
	return
}

// DownloadInDirectory downloads in parallel a set of src paths and saves them in a dst directory
func (d *Downloader) DownloadInDirectory(ctx context.Context, dst string, paths ...string) error {
	return d.Download(ctx, paths, func(ctx context.Context, idx int, path string, r io.ReadCloser) (err error) {
		// Make sure to close the reader
		defer r.Close()

		// Make sure destination directory exists
		if err = os.MkdirAll(dst, 0700); err != nil {
			err = errors.Wrapf(err, "astihttp: mkdirall %s failed", dst)
			return
		}

		// Create destination file
		var f *os.File
		dst := filepath.Join(dst, filepath.Base(path))
		if f, err = os.Create(dst); err != nil {
			err = errors.Wrapf(err, "astihttp: creating %s failed", dst)
			return
		}
		defer f.Close()

		// Copy
		if _, err = astiio.Copy(ctx, r, f); err != nil {
			err = errors.Wrapf(err, "astihttp: copying content to %s failed", dst)
			return
		}
		return
	})
}

type chunk struct {
	idx  int
	r    io.ReadCloser
	path string
}

// DownloadInWriter downloads in parallel a set of src paths and concatenates them in order in a writer
func (d *Downloader) DownloadInWriter(ctx context.Context, w io.Writer, paths ...string) (err error) {
	// Download
	var cs []chunk
	var m sync.Mutex // Locks cs
	var requiredIdx int
	err = d.Download(ctx, paths, func(ctx context.Context, idx int, path string, r io.ReadCloser) (err error) {
		// Lock
		m.Lock()
		defer m.Unlock()

		// Check where to insert chunk
		var idxInsert = -1
		for idxChunk := 0; idxChunk < len(cs); idxChunk++ {
			if idx < cs[idxChunk].idx {
				idxInsert = idxChunk
				break
			}
		}

		// Create chunk
		c := chunk{
			idx:  idx,
			path: path,
			r:    r,
		}

		// Add chunk
		if idxInsert > -1 {
			cs = append(cs[:idxInsert], append([]chunk{c}, cs[idxInsert:]...)...)
		} else {
			cs = append(cs, c)
		}

		// Loop through chunks
		for idxChunk := 0; idxChunk < len(cs); idxChunk++ {
			// Get chunk
			c := cs[idxChunk]

			// The chunk should be copied
			if c.idx == requiredIdx {
				// Copy chunk content
				_, err = astiio.Copy(ctx, c.r, w)

				// Make sure the reader is closed
				c.r.Close()

				// Remove chunk
				requiredIdx++
				cs = append(cs[:idxChunk], cs[idxChunk+1:]...)
				idxChunk--

				// Check error now so that chunk is still removed and reader is closed
				if err != nil {
					err = errors.Wrapf(err, "astihttp: copying chunk #%d to dst failed", c.idx)
					return
				}
			}
		}
		return
	})

	// Make sure to close all readers
	for _, c := range cs {
		c.r.Close()
	}
	return
}

// DownloadInFile downloads in parallel a set of src paths and concatenates them in order in a writer
func (d *Downloader) DownloadInFile(ctx context.Context, dst string, paths ...string) (err error) {
	// Make sure destination directory exists
	if err = os.MkdirAll(filepath.Dir(dst), 0700); err != nil {
		err = errors.Wrapf(err, "astihttp: mkdirall %s failed", filepath.Dir(dst))
		return
	}

	// Create destination file
	var f *os.File
	if f, err = os.Create(dst); err != nil {
		err = errors.Wrapf(err, "astihttp: creating %s failed", dst)
		return
	}
	defer f.Close()

	// Download in writer
	return d.DownloadInWriter(ctx, f, paths...)
}
