package astipcm

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAudioLevel(t *testing.T) {
	assert.Equal(t, 2.160246899469287, AudioLevel([]int{1, 2, 3}))
}
