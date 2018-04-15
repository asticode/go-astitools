package astiunicode

import (
	"unicode/utf8"
)

// StripUTF8Chars strips UTF8 chars
func StripUTF8Chars(i []byte) (o []byte) {
	buf := i
	for len(buf) > 0 {
		r, size := utf8.DecodeRune(buf)
		if r == utf8.RuneError && size == 1 {
			buf = buf[size:]
			continue
		}
		o = append(o, buf[:size]...)
		buf = buf[size:]
	}
	return
}
