package hw09structvalidator

import "fmt"

type LenRule struct {
	Value int
}

func (r LenRule) Validate(v interface{}) error {
	if value, ok := v.(int); ok {
		v = fmt.Sprintf("%d", value)
	}

	if value, ok := v.(string); ok {
		if len(value) != r.Value {
			return fmt.Errorf("%w: expected %d, got %d", ErrStringLengthMismatch, r.Value, len(value))
		}
		return nil
	}
	return fmt.Errorf("%w: %T", ErrSysUnsupportedType, v)
}
