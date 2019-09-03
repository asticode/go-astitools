package astitemplate

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/pkg/errors"
)

// Templater represents an object capable of storing templates
type Templater struct {
	layouts   []string
	m         sync.Mutex
	templates map[string]*template.Template
}

// NewTemplater creates a new templater
func NewTemplater(templatesPath, layoutsPath, ext string) (t *Templater, err error) {
	// Create templater
	t = &Templater{templates: make(map[string]*template.Template)}

	// Get layouts
	if err = filepath.Walk(layoutsPath, func(path string, info os.FileInfo, e error) (err error) {
		// Check input error
		if e != nil {
			err = errors.Wrapf(e, "astitemplate: walking layouts has an input error for path %s", path)
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
		t.layouts = append(t.layouts, path)
		return
	}); err != nil {
		err = errors.Wrapf(err, "astitemplate: walking layouts in %s failed", layoutsPath)
		return
	}

	// Loop through templates
	if err = filepath.Walk(templatesPath, func(path string, info os.FileInfo, e error) (err error) {
		// Check input error
		if e != nil {
			err = errors.Wrapf(e, "astitemplate: walking templates has an input error for path %s", path)
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
			err = errors.Wrapf(err, "astitemplate: reading template content of %s failed", path)
			return
		}

		// Add template
		// We use ToSlash to homogenize Windows path
		if err = t.Add(filepath.ToSlash(strings.TrimPrefix(path, templatesPath)), string(b)); err != nil {
			err = errors.Wrap(err, "astitemplate: adding template failed")
			return
		}
		return
	}); err != nil {
		err = errors.Wrapf(err, "astitemplate: walking templates in %s failed", templatesPath)
		return
	}
	return
}

// Add adds a new template
func (t *Templater) Add(path, content string) (err error) {
	// Parse
	var tpl *template.Template
	if tpl, err = t.Parse(content); err != nil {
		err = errors.Wrapf(err, "astitemplate: parsing template for path %s failed", path)
		return
	}

	// Add template
	t.m.Lock()
	t.templates[path] = tpl
	t.m.Unlock()
	return
}

// Del deletes a template
func (t *Templater) Del(path string) {
	t.m.Lock()
	defer t.m.Unlock()
	delete(t.templates, path)
}

// Template retrieves a templates
func (t *Templater) Template(path string) (tpl *template.Template, ok bool) {
	t.m.Lock()
	defer t.m.Unlock()
	tpl, ok = t.templates[path]
	return
}

func (t *Templater) Parse(content string) (o *template.Template, err error) {
	// Parse content
	o = template.New("root")
	if o, err = o.Parse(content); err != nil {
		err = errors.Wrap(err, "astitemplate: parsing template content failed")
		return
	}

	// Parse files
	if o, err = o.ParseFiles(t.layouts...); err != nil {
		err = errors.Wrapf(err, "astitemplate: parsing layouts %s failed", strings.Join(t.layouts, ", "))
		return
	}
	return
}
