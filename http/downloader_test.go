package astihttp

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloader(t *testing.T) {
	// Init
	m := &sync.Mutex{}
	m.Lock()
	s := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case "/1":
			rw.Write([]byte("1"))
		case "/2":
			m.Lock()
			rw.Write([]byte("2"))
		case "/3":
			rw.Write([]byte("3"))
			m.Unlock()
		case "/4":
			rw.Write([]byte("4"))
		}
	}))
	defer s.Close()

	// Download in writer
	buf := &bytes.Buffer{}
	d := NewDownloader(DownloaderOptions{NumberOfWorkers: 2})
	err := d.DownloadInWriter(context.Background(), buf, s.URL+"/1", s.URL+"/2", s.URL+"/3", s.URL+"/4")
	assert.NoError(t, err)
	assert.Equal(t, "1234", buf.String())
}
