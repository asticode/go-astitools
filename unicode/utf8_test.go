package astiunicode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanUTF8Chars(t *testing.T) {
	assert.Equal(t, []byte("az"), CleanUTF8Chars([]byte("az")))
	assert.Equal(t, []byte("az"), CleanUTF8Chars([]byte("a\xc5z")))
}
