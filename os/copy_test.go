package astios

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	// Create temporary dir
	p, err := ioutil.TempDir("", "astitools_os_test_")
	if err != nil {
		t.Log(errors.Wrapf(err, "creating %s failed, skipping TestCopy", p))
		return
	}

	// Make sure the dir is deleted
	defer func() {
		return
		if err = os.RemoveAll(p); err != nil {
			t.Log(errors.Wrapf(err, "removing %s failed", p))
		}
	}()

	// Copy file
	err = Copy(context.Background(), "./testdata/copy/f", filepath.Join(p, "f"))
	assert.NoError(t, err)
	checkFile(t, filepath.Join(p, "f"), []byte("0"))

	// Copy dir
	err = Copy(context.Background(), "./testdata/copy/d", filepath.Join(p, "d"))
	assert.NoError(t, err)
	checkFile(t, filepath.Join(p, "d", "f1"), []byte("1"))
	checkFile(t, filepath.Join(p, "d", "d1", "f11"), []byte("2"))
	checkFile(t, filepath.Join(p, "d", "d2", "f21"), []byte("3"))
	checkFile(t, filepath.Join(p, "d", "d2", "d21", "f211"), []byte("4"))
}

func checkFile(t *testing.T, p string, c []byte) {
	b, err := ioutil.ReadFile(p)
	assert.NoError(t, err)
	assert.Equal(t, c, b)
}
