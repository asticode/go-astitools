package astierror

import "strings"

// Multiple is an object containing multiple errors
type Multiple []error

// NewMultiple creates new multiple errors
func NewMultiple(errs []error) Multiple {
	return Multiple(errs)
}

// Error implements the error interface
func (m Multiple) Error() string {
	var ss []string
	for _, err := range m {
		ss = append(ss, err.Error())
	}
	return strings.Join(ss, " | ")
}
