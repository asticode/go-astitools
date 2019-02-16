package astissh

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/asticode/go-astitools/defer"
	"github.com/asticode/go-astitools/io"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

// SessionFunc represents a function that can create a new ssh session
type SessionFunc func() (*ssh.Session, *astidefer.Closer, error)

// Copy is a cancellable copy
// If src is a file, dst must be the full path to file once copied
// If src is a dir, dst must be the full path to the dir once copied
func Copy(ctx context.Context, src, dst string, fn SessionFunc) (err error) {
	// Check context
	if err = ctx.Err(); err != nil {
		return
	}

	// Stat src
	var statSrc os.FileInfo
	if statSrc, err = os.Stat(src); err != nil {
		err = errors.Wrapf(err, "stating %s failed", src)
		return
	}

	// Dir
	if statSrc.IsDir() {
		if err = filepath.Walk(src, func(path string, info os.FileInfo, errWalk error) (err error) {
			// Check error
			if errWalk != nil {
				err = errWalk
				return
			}

			// Do not process root
			if src == path {
				return
			}

			// Copy
			var p = filepath.Join(dst, strings.TrimPrefix(path, filepath.Clean(src)))
			if err = Copy(ctx, path, p, fn); err != nil {
				err = errors.Wrapf(err, "copying %s to %s failed", path, p)
				return

			}
			return
		}); err != nil {
			return
		}
		return
	}

	// Create ssh session
	var s *ssh.Session
	var c *astidefer.Closer
	if s, c, err = fn(); err != nil {
		err = errors.Wrap(err, "main: creating ssh session failed")
		return
	}
	defer c.Close()

	// Create the destination folder
	if err = s.Run("mkdir -p " + filepath.Dir(dst)); err != nil {
		err = errors.Wrapf(err, "astissh: creating %s failed", filepath.Dir(dst))
		return
	}

	// Open file
	var f *os.File
	if f, err = os.Open(src); err != nil {
		err = errors.Wrapf(err, "astissh: opening %s failed", src)
		return
	}
	defer f.Close()

	// Create ssh session
	if s, c, err = fn(); err != nil {
		err = errors.Wrap(err, "main: creating ssh session failed")
		return
	}
	defer c.Close()

	// Create stdin pipe
	var stdin io.WriteCloser
	if stdin, err = s.StdinPipe(); err != nil {
		err = errors.Wrap(err, "astissh: creating stdin pipe failed")
		return
	}
	defer stdin.Close()

	// Run "cat" command
	if err = s.Start("cat > " + dst); err != nil {
		err = errors.Wrapf(err, "astissh: cat to %s failed", dst)
		return
	}

	// Copy
	if _, err = astiio.Copy(ctx, f, stdin); err != nil {
		err = errors.Wrap(err, "astissh: copying failed")
		return
	}

	// Cat waits for the newline symbol from stdin to perform writing
	if _, err = stdin.Write([]byte("\n")); err != nil {
		err = errors.Wrap(err, "astissh: adding newline symbol failed")
		return
	}

	// Create ssh session
	if s, c, err = fn(); err != nil {
		err = errors.Wrap(err, "main: creating ssh session failed")
		return
	}
	defer c.Close()

	// Remove last char
	if err = s.Run("truncate --size=-1 " + dst); err != nil {
		err = errors.Wrap(err, "astissh: removing last char failed")
		return
	}
	return
}
