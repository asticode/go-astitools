package astitemplate

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
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

// ParseHTML parses html templates
// TODO Handle recursive glob
func ParseHTML(i string) (ts map[string]*template.Template, err error) {
	// Init
	i = filepath.Clean(i)
	ts = make(map[string]*template.Template)

	// Get layouts
	var ls []string
	if ls, err = filepath.Glob(i + "/layouts/*.html"); err != nil {
		return
	}

	// Get pages
	var ps []string
	if ps, err = filepath.Glob(i + "/pages/*.html"); err != nil {
		return
	}

	// Build root template
	tr := template.New("root")
	if tr, err = tr.Parse(`{{ template "base" . }}`); err != nil {
		return
	}

	// Loop through pages
	for _, p := range ps {
		// Clone root template
		var t *template.Template
		if t, err = tr.Clone(); err != nil {
			return
		}

		// Parse files
		if t, err = t.ParseFiles(append(ls, p)...); err != nil {
			return
		}

		// Add template
		ts[strings.TrimPrefix(p, i)] = t
	}
	return
}
