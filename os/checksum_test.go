package astios

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChecksum(t *testing.T) {
	c, err := Checksum("./testdata/checksum")
	assert.NoError(t, err)
	assert.Equal(t, "cRDtpNCeBiql5KOQsKVyrA0sAiA=", c)
}
