package envalid

import (
	"fmt"
	"os"
)

/// Variable

// Variable describes the environment variable for validation.
type Variable[T any] struct {
	Key         string // Env key (e.g. TOKEN)
	Description string // Optional description used in errors
	Default     *T     // Optional default value returned if env is not set
}

// WithDefault provides a default value for a Variable.
func WithDefault[T any](v T) *T {
	return &v
}

/// Validator

// Validator read Parser.
type Validator[T any] = Parser[T]

// Str read ParseStr.
var Str = ParseStr

// Int32 read ParseInt32.
var Int32 = ParseInt32

// Int64 read ParseInt64.
var Int64 = ParseInt64

// Float32 read ParseFloat32.
var Float32 = ParseFloat32

// Float64 read ParseFloat64.
var Float64 = ParseFloat64

// Bool read ParseBool.
var Bool = ParseBool

// Port read ParsePort.
var Port = ParsePort

/// Error

type ErrorCode int

const (
	_ ErrorCode = iota
	// ErrNoValue occurs when env is not set and no default value is specified.
	ErrNoValue
	// ErrInvalidValue occurs when env is invalid for the target type.
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
// It returns the default value if env is not set and a default value is specified;
// a ValidationError with ErrNoValue if env is not set and no default value is specified;
// a ValidationError with ErrInvalidValue if env is invalid for the target type.
func Validate[T any](validator Validator[T], variable Variable[T]) (T, error) {
	var zero T

	env, exists := os.LookupEnv(variable.Key)
	if !exists {
		if variable.Default != nil {
			return *variable.Default, nil
		}

		return zero, newValidationError(ErrNoValue, variable, nil)
	}

	val, err := validator(env)
	if err != nil {
		return zero, newValidationError(ErrInvalidValue, variable, err)
	}

	return val, nil
}
