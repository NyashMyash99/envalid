package envalid

import (
	"errors"
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
	"unsafe"
)

// Parser describes a generic function that parses a string into a value of type T.
//
// It returns an error if the input is invalid for the target type.
type Parser[T any] func(string) (T, error)

// TrimmedParser wraps the Parser, preprocessing the input with the strings.TrimSpace function.
func TrimmedParser[T any](p Parser[T]) Parser[T] {
	return func(s string) (T, error) {
		return p(strings.TrimSpace(s))
	}
}

// ParseStr parses string s as a trimmed string.
var ParseStr = TrimmedParser(func(s string) (string, error) {
	return s, nil
})

/// ParseNumber

func ternary[T any](cond bool, t T, f T) T {
	if cond {
		return t
	}
	return f
}

// mapStrconvErr converts strconv errors into more detailed library errors.
func mapStrconvErr(err error, typ string, minVal, maxVal any) error {
	if errors.Is(err, strconv.ErrSyntax) {
		return fmt.Errorf("value must be a %s", typ)
	}
	if errors.Is(err, strconv.ErrRange) {
		return fmt.Errorf("value must be a number between %v and %v", minVal, maxVal)
	}
	return err
}

func parseIntGeneric[T ~int32 | ~int64](s string) (T, error) {
	var zero T
	bitSize := int(unsafe.Sizeof(zero) * 8)

	v, err := strconv.ParseInt(s, 10, bitSize)
	if err != nil {
		maxVal := ternary(bitSize == 32, math.MaxInt32, math.MaxInt64)
		return 0, mapStrconvErr(err, fmt.Sprintf("int%d", bitSize), -maxVal-1, maxVal)
	}
	return T(v), nil
}

func parseFloatGeneric[T ~float32 | ~float64](s string) (T, error) {
	var zero T
	bitSize := int(unsafe.Sizeof(zero) * 8)

	v, err := strconv.ParseFloat(s, bitSize)
	if err != nil {
		maxVal := ternary(bitSize == 32, math.MaxFloat32, math.MaxFloat64)
		return 0, mapStrconvErr(err, fmt.Sprintf("float%d", bitSize), -maxVal, maxVal)
	}
	return T(v), nil
}

// ParseInt32 parses the string s as an integer.
//
// It returns an error if the string s is not a number or out of range.
var ParseInt32 = TrimmedParser(func(s string) (int32, error) {
	v, err := parseIntGeneric[int32](s)
	return v, err
})

// ParseInt64 parses the string s as an integer.
//
// It returns an error if the string s is not a number or out of range.
var ParseInt64 = TrimmedParser(func(s string) (int64, error) {
	v, err := parseIntGeneric[int64](s)
	return v, err
})

// ParseFloat32 parses the string s as a float.
//
// It returns an error if the string s is not a number or out of range.
var ParseFloat32 = TrimmedParser(func(s string) (float32, error) {
	v, err := parseFloatGeneric[float32](s)
	return v, err
})

// ParseFloat64 parses the string s as a float.
//
// It returns an error if the string s is not a number or out of range.
var ParseFloat64 = TrimmedParser(func(s string) (float64, error) {
	v, err := parseFloatGeneric[float64](s)
	return v, err
})

/// ParseBool

// BoolMap maps bool-like strings to bool values.
var BoolMap = map[string]bool{
	"1": true, "true": true, "t": true, "yes": true, "y": true, "on": true,
	"0": false, "false": false, "f": false, "no": false, "n": false, "off": false,
}

// ParseBool parses the string s as a bool.
//
// It returns an error if the string s is not a bool.
var ParseBool = TrimmedParser(func(s string) (bool, error) {
	s = strings.ToLower(s)
	if v, ok := BoolMap[s]; ok {
		return v, nil
	}
	return false, errors.New("value must be a bool-like")
})

/// ParsePort

const (
	minPort = 1
	maxPort = 65535
)

// ParsePort parses the string s as a TCP/UDP port (1–65535).
//
// It returns an error if the string s is not a number or out of range.
var ParsePort = TrimmedParser(func(s string) (uint16, error) {
	v, err := parseIntGeneric[int32](s)
	if err != nil || v < minPort || v > maxPort {
		return 0, fmt.Errorf("port must be a number between %d and %d", minPort, maxPort)
	}
	return uint16(v), nil
})

///

// ParseIPv4 parses the string s as a IPv4 address.
//
// It returns an error if the string s is not a valid IPv4 address.
var ParseIPv4 = TrimmedParser(func(s string) (net.IP, error) {
	v := net.ParseIP(s)
	if v == nil || v.To4() == nil {
		return nil, errors.New("value must be a valid IPv4 address")
	}
	return v, nil
})

// ParseIPv6 parses the string s as a IPv6 address.
//
// It returns an error if the string s is not a valid IPv6 address.
var ParseIPv6 = TrimmedParser(func(s string) (net.IP, error) {
	v := net.ParseIP(s)
	if v == nil || v.To16() == nil || v.To4() != nil {
		return nil, errors.New("value must be a valid IPv6 address")
	}
	return v, nil
})
}

var _ Parser[uint16] = ParsePort
