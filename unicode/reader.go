package astiunicode

import (
	"bufio"
	"bytes"
	"io"

	"github.com/pkg/errors"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// BOM headers
var (
	bomHeaderUTF32BE = []byte{0x00, 0x00, 0xFE, 0xFF}
	bomHeaderUTF32LE = []byte{0xFE, 0xFF, 0x00, 0x00}
	bomHeaderUTF16BE = []byte{0xFE, 0xFF}
	bomHeaderUTF16LE = []byte{0xFF, 0xFE}
	bomHeaderUTF8    = []byte{0xEF, 0xBB, 0xBF}
)

// NewReader creates a new unicode reader
func NewReader(i io.Reader) (o io.Reader, err error) {
	// Create reader
	r := bufio.NewReader(i)

	// Read first 4 bytes
	var b []byte
	if b, err = r.Peek(4); err != nil {
		err = errors.Wrap(err, "astiunicode: reading first 4 bytes failed")
		return
	}

	// Create transformer
	if bytes.HasPrefix(b, bomHeaderUTF32BE) || bytes.HasPrefix(b, bomHeaderUTF32LE) {
		err = errors.New("astiunicode: UTF32 is not handled yet")
		return
	} else if bytes.HasPrefix(b, bomHeaderUTF16BE) {
		o = transform.NewReader(r, unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder())
	} else if bytes.HasPrefix(b, bomHeaderUTF16LE) {
		o = transform.NewReader(r, unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder())
	} else if bytes.HasPrefix(b, bomHeaderUTF8) {
		o = transform.NewReader(r, unicode.UTF8.NewDecoder())
	} else {
		o = r
	}
	return
}
