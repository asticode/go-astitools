package astifloat

import (
	"bytes"
	"strconv"

	"github.com/pkg/errors"
)

// Rational represents a rational
type Rational struct{ Den, Num int }

// UnmarshalText implements the TextUnmarshaler interface
func (r *Rational) UnmarshalText(b []byte) (err error) {
	r.Num = 0
	r.Den = 1
	items := bytes.Split(b, []byte("/"))
	if r.Num, err = strconv.Atoi(string(items[0])); err != nil {
		err = errors.Wrapf(err, "astiencoder: atoi of %s failed", string(items[0]))
		return
	}
	if len(items) > 1 {
		if r.Den, err = strconv.Atoi(string(items[1])); err != nil {
			err = errors.Wrapf(err, "astiencoder: atoi of %s failed", string(items[1]))
			return
		}
	}
	return
}
