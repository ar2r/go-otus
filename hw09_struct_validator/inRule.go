package hw09structvalidator

import (
	"fmt"
	"strconv"
)

type InRule struct {
	Values []string
}

func (r InRule) Validate(v interface{}) error {
	if value, ok := v.(int); ok {
		v = strconv.Itoa(value)
	}
	vAsString, ok := v.(string)
	if !ok {
		return fmt.Errorf("%w: %T", ErrSysUnsupportedType, v)
	}

	isMatched := false

	for _, item := range r.Values {
		if item == vAsString {
			isMatched = true
		}
	}
	if !isMatched {
		return fmt.Errorf("%w: %s in %v", ErrValueNotInList, vAsString, r.Values)
	}
	return nil
}
