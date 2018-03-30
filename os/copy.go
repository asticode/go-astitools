package astios

import (
	"context"
	"os"

	"path/filepath"

	"github.com/asticode/go-astitools/io"
)

// Copy is a cross partitions cancellable copy
func Copy(ctx context.Context, src, dst string) (err error) {
	var (
		srcFile  *os.File
		dstFile  *os.File
		srcStats os.FileInfo
	)

	// Check context
	if err = ctx.Err(); err != nil {
		return
	}

	// Open the source file
	if srcFile, err = os.Open(src); err != nil {
		return
	}
	defer srcFile.Close()

	// Get the file stats/mode
	if srcStats, err = srcFile.Stat(); err != nil {
		return err
	}

	// Check context
	if err = ctx.Err(); err != nil {
		return
	}

	// Create the destination folder
	if err = os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return
	}

	// Check context
	if err = ctx.Err(); err != nil {
		return
	}

	// Create the destination file
	if dstFile, err = os.Create(dst); err != nil {
		return
	}
	defer dstFile.Close()

	if err = dstFile.Chmod(srcStats.Mode()); err != nil {
		return
	}

	// Check context
	if err = ctx.Err(); err != nil {
		return
	}

	// Copy the content
	if _, err = astiio.Copy(ctx, srcFile, dstFile); err != nil {
		return
	}

	return
}
