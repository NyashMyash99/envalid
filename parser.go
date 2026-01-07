package envalid

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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
func trim(s string) string {
	return strings.TrimSpace(s)
}

func ParseStr(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}

	return s, nil
}

var _ Parser[string] = ParseStr

// ParsePort
/// ParseNumber

func mapStrconvErr(err error, typ string) error {
	if errors.Is(err, strconv.ErrSyntax) {
		return fmt.Errorf("value must be a %s", typ)
	}
	if errors.Is(err, strconv.ErrRange) {
		return fmt.Errorf("value must be in range of %s", typ)
	}
	return err
}

func parseIntGeneric(s string, bitSize int) (int64, error) {
	v, err := strconv.ParseInt(trim(s), 10, bitSize)
	if err != nil {
		return 0, mapStrconvErr(err, fmt.Sprintf("int%d", bitSize))
	}
	return v, nil
}

const (
	minPort = 1
	maxPort = 65535
)

// ParsePort parses string s as a TCP/UDP port (1–65535).
//
// It returns an error if string s is not a number or out of range.
func ParsePort(s string) (uint16, error) {
	v, err := parseIntGeneric(s, 16)
	if err != nil || v < minPort || v > maxPort {
		return 0, fmt.Errorf("port must be a number between %d and %d", minPort, maxPort)
	}
	return uint16(v), nil
}

var _ Parser[uint16] = ParsePort
