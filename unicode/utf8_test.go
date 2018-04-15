package astiunicode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripUTF8Chars(t *testing.T) {
	assert.Equal(t, []byte("az"), StripUTF8Chars([]byte("az")))
	assert.Equal(t, []byte("az"), StripUTF8Chars([]byte("a\xc5z")))
}
