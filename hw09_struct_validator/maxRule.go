package hw09structvalidator

import "fmt"

type MaxRule struct {
	Value int
}

func (r MaxRule) Validate(v interface{}) error {
	if value, ok := v.(int); ok {
		if value > r.Value {
			return fmt.Errorf("%w: %d", ErrValueIsMoreThanMaxValue, r.Value)
		}
		return nil
	}
	return fmt.Errorf("%w: %T", ErrSysUnsupportedType, v)
}
