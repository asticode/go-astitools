package astislice_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/asticode/go-astitools/slice"
)

func TestInStringSlice(t *testing.T) {
	assert.False(t, astislice.InStringSlice("test", []string{"test1", "test2"}))
	assert.True(t, astislice.InStringSlice("test1", []string{"test1", "test2"}))
}
