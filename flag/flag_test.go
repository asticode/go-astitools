package astiflag_test

import (
	"os"
	"testing"

	"github.com/asticode/go-toolkit/flag"
	"github.com/stretchr/testify/assert"
)

func TestSubcommand(t *testing.T) {
	os.Args = []string{"bite"}
	assert.Equal(t, "", flag.Subcommand())
	os.Args = []string{"bite", "-caca"}
	assert.Equal(t, "", flag.Subcommand())
	os.Args = []string{"bite", "caca"}
	assert.Equal(t, "caca", flag.Subcommand())
}
