package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/exp/constraints"
)

var (
	ErrSysNotAStruct          = SystemError{errors.New("not a struct")}
	ErrSysUnsupportedType     = SystemError{errors.New("unsupported type")}
	ErrSysUnsupportedSlice    = SystemError{errors.New("slice of unsupported type")}
	ErrSysInvalidRule         = SystemError{errors.New("invalid rule")}
	ErrSysCantConvertLenValue = SystemError{errors.New("can't convert len value")}
	ErrSysCantConvertMaxValue = SystemError{errors.New("can't convert max value")}
	ErrSysCantConvertMinValue = SystemError{errors.New("can't convert min value")}
	ErrSysRegexpCompile       = SystemError{errors.New("regexp compile failed")}
)

var (
	ErrValueIsLessThanMinValue = errors.New("value is less than min value")
	ErrValueIsMoreThanMaxValue = errors.New("value is more than max value")
	ErrValueNotInList          = errors.New("value is not in the list")
	ErrStringLengthMismatch    = errors.New("string length mismatch")
	ErrRegexpMatchFailed       = errors.New("regexp match failed")
)

var (
	MinRuleTag    = "min"
	MaxRuleTag    = "max"
	InRuleTag     = "in"
	LenRuleTag    = "len"
	RegexpRuleTag = "regexp"
)

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Integer interface {
	Signed | Unsigned
}

type SystemError struct {
	Err error
}

func (e SystemError) Error() string {
	return fmt.Sprintf("system error: %v", e.Err)
}

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	builder := strings.Builder{}
	for _, e := range v {
		builder.WriteString(e.Field + ": " + e.Err.Error() + "\n")
	}
	return builder.String()
}

func Validate(v interface{}) error {
	errorsSlice := make(ValidationErrors, 0)

	vType := reflect.TypeOf(v)
	if vType.Kind() != reflect.Struct {
		return fmt.Errorf("%w: %v", ErrSysNotAStruct, vType.Kind())
	}

	for i := 0; i < vType.NumField(); i++ {
		propType := vType.Field(i)
		propValue := reflect.ValueOf(v).Field(i)
		propTagValue := propType.Tag.Get("validate")

		if propTagValue == "" {
			continue
		}

		rules, parseErr := parseRules(propTagValue)
		if parseErr != nil {
			return parseErr
		}

		var errorsStack []error
		var err error

		//nolint:exhaustive
		switch propValue.Kind() {
		case reflect.String:
			errorsStack, err = validateValues([]string{propValue.String()}, rules)
		case reflect.Int:
			errorsStack, err = validateValues([]int{int(propValue.Int())}, rules)
		case reflect.Slice:
			switch propValue.Type().Elem().Kind() {
			case reflect.String:
				errorsStack, err = validateValues(propValue.Interface().([]string), rules)
			case reflect.Int:
				errorsStack, err = validateValues(propValue.Interface().([]int), rules)
			default:
				return fmt.Errorf("%w: %v", ErrSysUnsupportedSlice, propValue.Type().Elem().Kind())
			}
		default:
			return fmt.Errorf("%w: %v", ErrSysUnsupportedType, propValue.Kind())
		}

		// Обработка ошибок и склеиваем их в один массив
		if err != nil {
			if errors.As(err, &SystemError{}) {
				return err
			}
		}
		for _, err := range errorsStack {
			errorsSlice = append(errorsSlice, ValidationError{
				Field: propType.Name,
				Err:   err,
			})
		}
	}

	if len(errorsSlice) > 0 {
		return errorsSlice
	}
	return nil
}

func validateValues[T interface{ constraints.Integer | string }](values []T, rules []Rule) ([]error, error) {
	var errorsSlice []error
	for _, value := range values {
		for _, rule := range rules {
			err := rule.Validate(value)
			if err != nil {
				if errors.As(err, &SystemError{}) {
					return nil, err
				}
				errorsSlice = append(errorsSlice, err)
			}
		}
	}

	if len(errorsSlice) > 0 {
		return errorsSlice, nil
	}

	return nil, nil
}
