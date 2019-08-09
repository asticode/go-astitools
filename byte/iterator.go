package astibyte

import "fmt"

// Iterator represents an object capable of iterating sequentially and safely through a slice of bytes
type Iterator struct {
	bs     []byte
	offset int
}

// NewIterator creates a new iterator
func NewIterator(bs []byte) *Iterator {
	return &Iterator{bs: bs}
}

// NextByte returns the slice's next byte
func (i *Iterator) NextByte() (b byte, err error) {
	if len(i.bs) < i.offset+1 {
		err = fmt.Errorf("astits: slice length is %d, offset %d is invalid", len(i.bs), i.offset)
		return
	}
	b = i.bs[i.offset]
	i.offset++
	return
}

// NextBytes returns the n slice's next bytes
func (i *Iterator) NextBytes(n int) (bs []byte, err error) {
	if len(i.bs) < i.offset+n {
		err = fmt.Errorf("astits: slice length is %d, offset %d is invalid", len(i.bs), i.offset+n)
		return
	}
	bs = make([]byte, n)
	copy(bs, i.bs[i.offset:i.offset+n])
	i.offset += n
	return
}

// Seek sets the iterator's offset
func (i *Iterator) Seek(offset int) {
	i.offset = offset
}

// FastForward increments the iterator's offset
func (i *Iterator) FastForward(delta int) {
	i.offset += delta
}

// HasBytesLeft checks whether the slice has some bytes left
func (i *Iterator) HasBytesLeft() bool {
	return i.offset < len(i.bs)
}

// Offset returns the offset
func (i *Iterator) Offset() int {
	return i.offset
}

// Dump dumps the rest of the slice
func (i *Iterator) Dump() (bs []byte) {
	if !i.HasBytesLeft() {
		return
	}
	bs = make([]byte, len(i.bs) - i.offset)
	copy(bs, i.bs[i.offset:len(i.bs)])
	i.offset = len(i.bs)
	return
}

// Len returns the slice length
func (i *Iterator) Len() int {
	return len(i.bs)
}
