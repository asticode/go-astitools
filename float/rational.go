package astifloat

import (
	"bytes"
	"strconv"

	"github.com/pkg/errors"
)

// Rational represents a rational
type Rational struct{ den, num int }

// NewRational creates a new rational
func NewRational(num, den int) *Rational {
	return &Rational{
		den: den,
		num: num,
	}
}

// Num returns the rational num
func (r *Rational) Num() int {
	return r.num
}

// Den returns the rational den
func (r *Rational) Den() int {
	return r.den
}

// ToFloat64 returns the rational as a float64
func (r *Rational) ToFloat64() float64 {
	return float64(r.num) / float64(r.den)
}

// UnmarshalText implements the TextUnmarshaler interface
func (r *Rational) UnmarshalText(b []byte) (err error) {
	r.num = 0
	r.den = 1
	if len(b) == 0 {
		return
	}
	items := bytes.Split(b, []byte("/"))
	if r.num, err = strconv.Atoi(string(items[0])); err != nil {
		err = errors.Wrapf(err, "astifloat: atoi of %s failed", string(items[0]))
		return
	}
	if len(items) > 1 {
		if r.den, err = strconv.Atoi(string(items[1])); err != nil {
			err = errors.Wrapf(err, "astifloat: atoi of %s failed", string(items[1]))
			return
		}
	}
	return
}
