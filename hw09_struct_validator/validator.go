package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
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

//nolint:gocognit
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

		switch propValue.Kind() {
		case reflect.String:
			errorsStack, err = stringValidate([]string{propValue.String()}, rules)
		case reflect.Int:
			errorsStack, err = intValidate([]int{int(propValue.Int())}, rules)
		case reflect.Slice:
			switch propValue.Type().Elem().Kind() {
			case reflect.String:
				errorsStack, err = stringValidate(propValue.Interface().([]string), rules)
			case reflect.Int:
				errorsStack, err = intValidate(propValue.Interface().([]int), rules)
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

func validate(values interface{}, rules []Rule) ([]error, error) {
	var errorsSlice []error

	// Кучу времени потратил, но так и не смог никак убрать это дублирование. Подскажите, пожалуйста, как это сделать
	switch v := values.(type) {
	case []int:
		for _, value := range v {
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
	case []string:
		for _, value := range v {
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
	default:
		return nil, fmt.Errorf("%w: %T", ErrSysUnsupportedType, values)
	}

	if len(errorsSlice) > 0 {
		return errorsSlice, nil
	}

	return nil, nil
}

func intValidate(values []int, rules []Rule) ([]error, error) {
	return validate(values, rules)
}

func stringValidate(values []string, rules []Rule) ([]error, error) {
	return validate(values, rules)
}
