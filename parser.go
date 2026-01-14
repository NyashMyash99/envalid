package envalid

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net"
	"net/mail"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"golang.org/x/net/idna"
	"golang.org/x/net/publicsuffix"
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

/// ParseDomain

func ptr[T any](v T) *T { return &v }

// PDomain describes a parsed domain.
type PDomain struct {
	String    string  // e.g. docs.nyashmyash99.dev
	Subdomain *string // docs (see String)
	Hostname  string  // nyashmyash99.dev (see String)
	Zone      *string // dev (see String)
}

// ParseDomain parses the string s as a domain (including IDN).
//
// It returns an error if the string s is not a valid domain.
var ParseDomain = TrimmedParser(func(s string) (*PDomain, error) {
	if strings.HasSuffix(s, ".") {
		return nil, errors.New("value must not contain a trailing dot")
	}

	if strings.IndexByte(s, '.') == -1 {
		return nil, errors.New("value must contain a zone")
	}

	v, err := idna.Registration.ToASCII(s)
	if err != nil {
		return nil, errors.New("value must be a valid domain")
	}

	lastDot := strings.LastIndexByte(v, '.')
	nextDot := strings.LastIndexByte(v[:lastDot], '.')
	var subdomainPtr *string

	hostname := v
	if nextDot != -1 {
		hostname = v[nextDot+1:]
		sub := v[:nextDot]
		subdomainPtr = &sub
	}

	zone, _ := publicsuffix.PublicSuffix(s)

	return &PDomain{
		String:    v,
		Subdomain: subdomainPtr,
		Hostname:  hostname,
		Zone:      ptr(zone),
	}, nil
})

/// ParseHost

// PHost describes a parsed host.
type PHost struct {
	String string   // e.g. docs.nyashmyash99.dev, 127.0.0.1 or ::1
	Domain *PDomain // Optional domain, depending on the host type
	IP     net.IP   // Optional ip, depending on the host type
}

// ParseHost parses a string as a host, which can be:
// - an IPv4 address;
// - an IPv6 address;
// - a domain (including localhost).
//
// It returns an error if the string s is not a host.
var ParseHost = TrimmedParser(func(s string) (*PHost, error) {
	if strings.EqualFold(s, "localhost") {
		d := &PDomain{String: "localhost", Hostname: "localhost"}
		return &PHost{String: "localhost", Domain: d}, nil
	}

	if v, err := ParseIPv4(s); err == nil {
		return &PHost{String: v.String(), IP: v}, nil
	}
	if v, err := ParseIPv6(s); err == nil {
		return &PHost{String: v.String(), IP: v}, nil
	}
	if v, err := ParseDomain(s); err == nil {
		return &PHost{String: v.String, Domain: v}, nil
	}
	return nil, errors.New("value must be a host (IPv4, IPv6, or domain, including localhost)")
})

/// ParseURL

// PURL describes a parsed url.
type PURL struct {
	String   string     // e.g. http://user:pass@docs.nyashmyash99.dev:443/envalid?tab=documentation#quick-start
	Protocol string     // http (see String)
	Username *string    // user (see String)
	Password *string    // pass (see String)
	Host     PHost      // docs.nyashmyash99.dev (see String)
	Port     *uint16    // 443 (see String)
	Path     *string    // envalid (see String)
	Params   url.Values // tab="documentation" (see String)
	Anchor   *string    // quick-start
}

// ParseURL parses the string s as a url in the maximum format "protocol://username:password@host:port/path?params#anchor".
//
// It returns an error if the string s is not a valid url.
var ParseURL = TrimmedParser(func(s string) (*PURL, error) {
	uri, err := url.Parse(s)
	if err != nil {
		return nil, errors.New("value must be a URL")
	}

	if uri.Scheme == "" {
		return nil, errors.New("value must be in absolute format (protocol://...)")
	}

	if uri.Opaque != "" {
		return nil, errors.New("value must be in default url format")
	}

	if h := uri.Hostname(); h == "" {
		return nil, errors.New("value must contain a host")
	}

	var portPtr *uint16
	if p := uri.Port(); p != "" {
		pp, err := ParsePort(p)
		if err != nil {
			return nil, err
		}
		portPtr = &pp
	}

	host, err := ParseHost(uri.Hostname())
	if err != nil {
		return nil, errors.New("hostname must be a valid host (IPv4, IPv6, or domain, including localhost)")
	}

	var usernamePtr, passwordPtr *string
	if uri.User != nil {
		usernamePtr = ptr(uri.User.Username())
		p, _ := uri.User.Password()
		passwordPtr = &p
	}

	var pathPtr *string
	if p := uri.Path; p != "" && p != "/" {
		pathPtr = ptr(strings.TrimPrefix(p, "/"))
	}

	var fragmentPtr *string
	if uri.Fragment != "" {
		fragmentPtr = ptr(uri.Fragment)
	}

	return &PURL{
		String:   s,
		Protocol: strings.ToLower(uri.Scheme),
		Username: usernamePtr,
		Password: passwordPtr,
		Host:     *host,
		Port:     portPtr,
		Path:     pathPtr,
		Params:   uri.Query(),
		Anchor:   fragmentPtr,
	}, nil
})

/// ParseJSON

// ParseJSON parses the string s as an JSON with schema T.
//
// It returns an error if the string s is not a valid JSON.
func ParseJSON[T any](s string) (T, error) {
	s = strings.TrimSpace(s)

	dec := json.NewDecoder(bytes.NewBufferString(s))
	dec.DisallowUnknownFields()

	var schema T
	if err := dec.Decode(&schema); err != nil {
		return schema, fmt.Errorf(
			"value must be an JSON with schema %q",
			reflect.TypeOf(schema).String(),
		)
	}

	return schema, nil
}

var _ Parser[any] = ParseJSON[any]

///

// ParseEmail parses the string s as an email address in the format "user@domain".
//
// It returns an error if the string s is not a valid email address.
var ParseEmail = TrimmedParser(func(s string) (string, error) {
	address, err := mail.ParseAddress(s)
	if err != nil {
		return "", errors.New("value must be an email")
	}

	// Cuts off:
	// - Daniil <contact@nyashmyash99.dev>
	// - "Daniil Koshkin"@nyashmyash99.dev
	hasDiff := address.Address != s
	// Cuts off "user@[127.0.0.1]".
	hasBracket := strings.IndexByte(address.Address, '[') != -1
	if hasDiff || hasBracket {
		return "", errors.New("value must be in format \"user@domain\"")
	}

	_, domain, _ := strings.Cut(address.Address, "@")
	if _, err := ParseDomain(domain); err != nil {
		return "", errors.New("email must contain a valid domain")
	}

	return s, nil
})
