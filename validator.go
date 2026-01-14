package envalid

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

/// Variable

// Variable describes the environment variable for validation.
//
// Options priority:
// 1. DevDefault (if GO_ENV != production)
// 2. Default
// 3. Variants
// 4. Validator
type Variable[T any] struct {
	Key         string   // Required env var key (e.g. GO_ENV)
	Description string   // Optional description used in errors
	Variants    []string // Optional array of available values for the env var. It is case-sensitive.
	Default     *T       // Optional default value returned if the env var is not set. See Variable for priority determination.
	DevDefault  *T       // Optional default value returned if the env var is not set, GO_ENV var is specified and is not "production". See Variable for priority determination.
}

// WithDefault provides a default value for a Variable.
func WithDefault[T any](v T) *T {
	return &v
}

/// Validator

// Validator see Parser.
type Validator[T any] = Parser[T]

// Str see ParseStr.
var Str = ParseStr

// Int32 see ParseInt32.
var Int32 = ParseInt32

// Int64 see ParseInt64.
var Int64 = ParseInt64

// Float32 see ParseFloat32.
var Float32 = ParseFloat32

// Float64 see ParseFloat64.
var Float64 = ParseFloat64

// Bool see ParseBool.
var Bool = ParseBool

// Port see ParsePort.
var Port = ParsePort

// IPv4 see ParseIPv4.
var IPv4 = ParseIPv4

// IPv6 see ParseIPv6.
var IPv6 = ParseIPv6

// Domain see ParseDomain.
var Domain = ParseDomain

// Host see ParseHost.
var Host = ParseHost

// URL see ParseURL.
var URL = ParseURL

// JSON see ParseJSON.
func JSON[T any](s string) (T, error) {
	return ParseJSON[T](s)
}

var _ Validator[any] = JSON[any]

// Email see ParseEmail.
var Email = ParseEmail

/// Error

type ErrorCode int

const (
	_ ErrorCode = iota
	// ErrNoValue occurs when an env var is not set and no default value is specified.
	ErrNoValue
	// ErrInvalidValue occurs when an env var is invalid for the target type.
	ErrInvalidValue
)

// ValidationError occurs when validation fails.
//
// It implements the error interface.
type ValidationError struct {
	Code ErrorCode
	msg  string
	err  error
}

func (e *ValidationError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.err)
	}
	return e.msg
}

func (e *ValidationError) Unwrap() error {
	return e.err
}

func newValidationError[T any](c ErrorCode, v Variable[T], err error) *ValidationError {
	msg := v.Key
	if v.Description != "" {
		msg += " (" + v.Description + ")"
	}
	return &ValidationError{
		Code: c,
		msg:  msg,
		err:  err,
	}
}

// Validate loads and validates an environment variable using the provided Validator.
//
// It returns the default value if the env var is not set and a default value is specified;
// the dev default value if the env var is not set, a dev default value is specified, GO_ENV var is specified and is not "production";
// a ValidationError with ErrNoValue if the env var is not set and no default value is specified;
// a ValidationError with ErrInvalidValue if the env var is not found in the available variants or is invalid for the target type.
func Validate[T any](validator Validator[T], variable Variable[T]) (T, error) {
	var zero T

	env, exists := os.LookupEnv(variable.Key)
	if !exists {
		goEnv, goEnvExists := os.LookupEnv("GO_ENV")
		if variable.DevDefault != nil && goEnvExists && !strings.EqualFold(goEnv, "production") {
			return *variable.DevDefault, nil
		}
		if variable.Default != nil {
			return *variable.Default, nil
		}
		return zero, newValidationError(ErrNoValue, variable, nil)
	}

	if len(variable.Variants) > 0 && !slices.Contains(variable.Variants, strings.TrimSpace(env)) {
		return zero, newValidationError(ErrInvalidValue, variable, nil)
	}

	val, err := validator(env)
	if err != nil {
		return zero, newValidationError(ErrInvalidValue, variable, err)
	}
	return val, nil
}
