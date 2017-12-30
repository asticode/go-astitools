package astitemplate

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

// ParseDirectory parses a directory recursively
func ParseDirectory(i, ext string) (t *template.Template, err error) {
	// Parse templates
	i = filepath.Clean(i)
	t = template.New("Root")
	return t, filepath.Walk(i, func(path string, info os.FileInfo, e error) (err error) {
		// Check input error
		if e != nil {
			err = e
			return
		}

		// Only process files
		if info.IsDir() {
			return
		}

		// Check extension
		if ext != "" && filepath.Ext(path) != ext {
			return
		}

		// Read file
		var b []byte
		if b, err = ioutil.ReadFile(path); err != nil {
			return
		}

		// Parse template
		var c = t.New(filepath.ToSlash(strings.TrimPrefix(path, i)))
		if _, err = c.Parse(string(b)); err != nil {
			return fmt.Errorf("%s while parsing template %s", err, path)
		}
		return
	})
}

// ParseDirectoryWithLayouts parses a directory recursively with layouts
func ParseDirectoryWithLayouts(templatesPath, layoutsPath, ext string) (ts map[string]*template.Template, err error) {
	// Init
	ts = make(map[string]*template.Template)

	// Get layouts
	var ls []string
	if err = filepath.Walk(layoutsPath, func(path string, info os.FileInfo, e error) (err error) {
		// Check input error
		if e != nil {
			err = errors.Wrapf(e, "walking layouts has an input error for path %s", path)
			return
		}

		// Only process files
		if info.IsDir() {
			return
		}

		// Check extension
		if ext != "" && filepath.Ext(path) != ext {
			return
		}

		// Append
		ls = append(ls, path)
		return
	}); err != nil {
		err = errors.Wrapf(err, "walking layouts in %s failed", layoutsPath)
		return
	}

	// Loop through templates
	if err = filepath.Walk(templatesPath, func(path string, info os.FileInfo, e error) (err error) {
		// Check input error
		if e != nil {
			err = errors.Wrapf(e, "walking templates has an input error for path %s", path)
			return
		}

		// Only process files
		if info.IsDir() {
			return
		}

		// Check extension
		if ext != "" && filepath.Ext(path) != ext {
			return
		}

		// Read file
		var b []byte
		if b, err = ioutil.ReadFile(path); err != nil {
			err = errors.Wrapf(err, "reading template content of %s failed", path)
			return
		}

		// Parse content
		var t = template.New("root")
		if t, err = t.Parse(string(b)); err != nil {
			err = errors.Wrapf(err, "parsing template content of %s failed", path)
			return
		}

		// Parse files
		if t, err = t.ParseFiles(ls...); err != nil {
			err = errors.Wrapf(err, "parsing layouts %s for %s failed", strings.Join(ls, ", "), path)
			return
		}

		// Add template
		ts[strings.TrimPrefix(path, templatesPath)] = t
		return
	}); err != nil {
		err = errors.Wrapf(err, "walking templates in %s failed", templatesPath)
		return
	}
	return
}
