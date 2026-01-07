package envalid

import (
	"fmt"
	"os"
)

// Variable describes an environment variable for validation.
type Variable[T any] struct {
	Key         string // Env key (e.g. TOKEN)
	Description string // Optional description used in errors
	Default     *T     // Optional default value returned if the env is not set
}

type ErrorCode int

const (
	_ ErrorCode = iota
	// ErrNoValue occurs when an env is not set and no default value is specified.
	ErrNoValue
	// ErrInvalidValue occurs when an env is invalid for the target type.
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

// Validate loads and validates an environment variable using the provided parser.
//
// It returns the default value if the env is not set and a default value is specified;
// a ValidationError with ErrNoValue if the env is not set and no default value is specified;
// a ValidationError with ErrInvalidValue if the env is invalid for the target type.
func Validate[T any](p Parser[T], v Variable[T]) (T, error) {
	var zero T

	env, exists := os.LookupEnv(v.Key)
	if !exists {
		if v.Default != nil {
			return *v.Default, nil
		}

		return zero, newValidationError(ErrNoValue, v, nil)
	}

	val, err := p(env)
	if err != nil {
		return zero, newValidationError(ErrInvalidValue, v, err)
	}

	return val, nil
}
