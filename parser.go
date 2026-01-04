package envalid

import (
	"errors"
	"strconv"
)

// ErrEmpty indicates that a value is empty.
var ErrEmpty = errors.New("empty value")

// ErrInvalid indicates that a value is invalid for the target type.
var ErrInvalid = errors.New("invalid value")

// Parser defines a generic function that parses a string into a value of type T.
//
// Returns ErrEmpty if value is empty, ErrInvalid if string is invalid for the target type.
type Parser[T any] func(string) (T, error)

// ParseStr

// ParseStr parses string s as a non-empty string.
//
// Returns ErrEmpty if string s is empty.
func ParseStr(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}

	return s, nil
}

var _ Parser[string] = ParseStr

// ParsePort

const (
	minPort = 1
	maxPort = 65535
)

// ParsePort parses string s as a port (1–65535).
//
// Returns ErrEmpty if string s is empty, ErrInvalid if string s is not a number or out of range.
func ParsePort(s string) (int, error) {
	if s == "" {
		return 0, ErrEmpty
	}

	i, err := strconv.Atoi(s)
	if err != nil || i < minPort || i > maxPort {
		return 0, ErrInvalid
	}

	return i, nil
}

var _ Parser[int] = ParsePort
