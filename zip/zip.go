package astizip

import (
	"archive/zip"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/asticode/go-astitools/io"
	"github.com/pkg/errors"
)

// Unzip unzips a src into a dst
// Possible src formats are /path/to/zip.zip or /path/to/zip.zip/internal/path if you only want to unzip files in
// /internal/path in the .zip archive
func Unzip(ctx context.Context, src, dst string) (err error) {
	// Parse src path
	var split = strings.Split(src, ".zip")
	src = split[0] + ".zip"
	var internalPath string
	if len(split) >= 2 {
		internalPath = split[1]
	}

	// Open overall reader
	var r *zip.ReadCloser
	if r, err = zip.OpenReader(src); err != nil {
		return errors.Wrapf(err, "opening overall zip reader on %s failed", src)
	}
	defer r.Close()

	// Loop through files
	for _, f := range r.File {
		// Validate internal path
		var n = string(os.PathSeparator) + f.Name
		if internalPath != "" && !strings.HasPrefix(n, internalPath) {
			continue
		}

		// Open file reader
		var fr io.ReadCloser
		if fr, err = f.Open(); err != nil {
			return errors.Wrapf(err, "opening zip reader on file %s failed", n)
		}
		defer fr.Close()

		// Update file path
		var p = filepath.Join(dst, strings.TrimPrefix(n, internalPath))
		if f.FileInfo().IsDir() {
			// If file is a dir we save its file mode
			if err = os.MkdirAll(p, f.FileInfo().Mode()); err != nil {
				return errors.Wrapf(err, "mkdirall %s failed", p)
			}
		} else {
			// Since dirs don't always come up we make sure the directory of the file exists
			if err = os.MkdirAll(filepath.Dir(p), 0775); err != nil {
				return errors.Wrapf(err, "mkdirall %s failed", filepath.Dir(p))
			}

			// Open the file
			var fl *os.File
			if fl, err = os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.FileInfo().Mode()); err != nil {
				return errors.Wrapf(err, "opening file %s failed", p)
			}
			defer fl.Close()

			// Copy
			if _, err = astiio.Copy(ctx, fr, fl); err != nil {
				return errors.Wrapf(err, "copying %s into %s failed", n, p)
			}
		}
	}
	return
}
