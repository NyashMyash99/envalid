package envalid

import (
	"errors"
)

// ErrEmpty indicates that a value is empty.
var ErrEmpty = errors.New("empty value")

// ParseStr parses string s as a non-empty string.
//
// Returns ErrEmpty if string s is empty.
func ParseStr(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}

	return s, nil
}
