package hw09structvalidator

import "fmt"

type MinRule struct {
	Value int
}

func (r MinRule) Validate(v interface{}) error {
	if value, ok := v.(int); ok {
		if value < r.Value {
			return fmt.Errorf("%w: %d", ErrValueIsLessThanMinValue, r.Value)
		}
		return nil
	}
	return fmt.Errorf("%w: %T", ErrSysUnsupportedType, v)
}
