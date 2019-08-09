package astibyte

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIterator(t *testing.T) {
	// Setup
	i := NewIterator([]byte("1234567"))

	// Length
	assert.Equal(t, 7, i.Len())

	// Next byte
	b, err := i.NextByte()
	assert.NoError(t, err)
	assert.Equal(t, byte('1'), b)

	// Next bytes
	bs, err := i.NextBytes(3)
	assert.NoError(t, err)
	assert.Equal(t, []byte("234"), bs)
	assert.Equal(t, 4, i.Offset())

	// Fast forward
	i.FastForward(2)
	assert.Equal(t, 6, i.Offset())
	assert.True(t, i.HasBytesLeft())

	// Last byte
	b, err = i.NextByte()
	assert.NoError(t, err)
	assert.Equal(t, byte('7'), b)
	assert.False(t, i.HasBytesLeft())

	// No bytes
	b, err = i.NextByte()
	assert.Error(t, err)

	// Dump
	i.FastForward(-2)
	assert.Equal(t, []byte("67"), i.Dump())

	// Seek
	i.Seek(2)
	b, err = i.NextByte()
	assert.NoError(t, err)
	assert.Equal(t, byte('3'), b)
	assert.True(t, i.HasBytesLeft())
}
