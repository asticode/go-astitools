package astitemplate_test

import (
	"testing"

	"github.com/asticode/go-toolkit/template"
	"github.com/stretchr/testify/assert"
)

func TestParseDirectory(t *testing.T) {
	tmpl, err := template.ParseDirectory("tests", "")
	assert.NoError(t, err)
	assert.Len(t, tmpl.Templates(), 4)
	tmpl, err = template.ParseDirectory("tests", ".html")
	assert.NoError(t, err)
	assert.Len(t, tmpl.Templates(), 3)
}
