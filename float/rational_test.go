package astifloat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRational(t *testing.T) {
	r := &Rational{}
	err := r.UnmarshalText([]byte(""))
	assert.Error(t, err)
	err = r.UnmarshalText([]byte("test"))
	assert.Error(t, err)
	err = r.UnmarshalText([]byte("1/test"))
	assert.Error(t, err)
	err = r.UnmarshalText([]byte("0"))
	assert.NoError(t, err)
	assert.Equal(t, 0, r.Num)
	assert.Equal(t, 1, r.Den)
	err = r.UnmarshalText([]byte("1/2"))
	assert.NoError(t, err)
	assert.Equal(t, 1, r.Num)
	assert.Equal(t, 2, r.Den)
}
