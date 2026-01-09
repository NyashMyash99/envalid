package envalid

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// ValidatorSupplier describes a function that provides a validator for the Schema.
type ValidatorSupplier func() (any, error)

// Supplier adapts the validator in ValidatorSupplier.
func Supplier[T any](validator Validator[T], variable Variable[T]) ValidatorSupplier {
	return func() (any, error) {
		return Validate(validator, variable)
	}
}

// Schema maps config field names to their validators.
type Schema map[string]ValidatorSupplier

// Load loads and validates environment variables, creating a config with type T using the provided Schema.
func Load[T any](s Schema) (*T, error) {
	var zero T
	rCfgType := reflect.TypeOf(zero)

	k := rCfgType.Kind()
	if k != reflect.Struct {
		return nil, fmt.Errorf("%s: got %s, want struct", rCfgType.Name(), k)
	}

	rCfg := reflect.New(rCfgType).Elem()
	var errs []error

	for i := 0; i < rCfgType.NumField(); i++ {
		rField := rCfgType.Field(i)
		fieldName := rField.Name

		validate := s[fieldName]
		if validate == nil {
			continue
		}

		val, err := validate()
		if err != nil {
			errs = append(errs, err)
			continue
		}

		rVal := reflect.ValueOf(val)
		if !rVal.Type().AssignableTo(rField.Type) {
			return nil, fmt.Errorf(
				"%s.%s: got %s, want %s",
				rCfgType.Name(),
				fieldName,
				rVal.Type(),
				rField.Type,
			)
		}

		rCfg.Field(i).Set(rVal)
	}

	var missingErrs, invalidErrs []error

	for _, err := range errs {
		var vErr *ValidationError
		if errors.As(err, &vErr) {
			switch vErr.Code {
			case ErrNoValue:
				missingErrs = append(missingErrs, err)
			case ErrInvalidValue:
				invalidErrs = append(invalidErrs, err)
			}
		}
	}

	var errMsg []string

	if len(missingErrs) > 0 {
		errMsg = append(errMsg, "Missing environment variables:")
		for _, err := range missingErrs {
			errMsg = append(errMsg, "\t"+err.Error())
		}
	}

	if len(invalidErrs) > 0 {
		errMsg = append(errMsg, "Invalid environment variables:")
		for _, err := range invalidErrs {
			errMsg = append(errMsg, "\t"+err.Error())
		}
	}

	if len(errMsg) > 0 {
		return nil, errors.New(strings.Join(errMsg, "\n"))
	}

	r := rCfg.Addr().Interface().(*T)
	return r, nil
}
