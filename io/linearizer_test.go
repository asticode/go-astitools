package astiio_test

import (
	"context"
	"io"
	"strconv"
	"testing"

	"github.com/asticode/go-astitools/io"
	"github.com/stretchr/testify/assert"
)

type linearizerReader struct {
	closed bool
	count  int
}

func (r *linearizerReader) Close() error {
	r.closed = true
	return nil
}

func (r *linearizerReader) Read(b []byte) (n int, err error) {
	if r.count == 3 {
		err = io.EOF
		return
	}
	b[0] = byte(strconv.Itoa(r.count)[0])
	b[1] = byte('t')
	b[2] = byte('e')
	b[3] = byte('s')
	b[4] = byte('t')
	b[5] = byte(strconv.Itoa(r.count)[0])
	n = 6
	r.count++
	return
}

func TestLinearizer(t *testing.T) {
	pr := &linearizerReader{}
	p := astiio.NewLinearizer(context.Background(), pr, 10, 30)
	go p.Start()
	b := make([]byte, 3)
	n, err := p.Read(b)
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, []byte("0te"), b)
	n, err = p.Read(b)
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, []byte("st0"), b)
	n, err = p.Read(b)
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, []byte("1te"), b)
	n, err = p.Read(b)
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, []byte("st1"), b)
	n, err = p.Read(b)
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, []byte("2te"), b)
	n, err = p.Read(b)
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, []byte("st2"), b)
	_, err = p.Read(b)
	assert.Error(t, err, io.EOF.Error())
}
