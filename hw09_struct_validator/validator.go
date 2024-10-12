package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
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

var (
	compiledRegexps = make(map[string]*regexp.Regexp)
	regexpMutex     = sync.Mutex{}
)

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

		rules, err := parseRules(propTagValue)
		if err != nil {
			return err
		}

		//nolint:exhaustive
		switch propValue.Kind() {
		case reflect.String:
			err := stringValidate(propValue.String(), propTagValue)
			if err != nil {
				if errors.As(err, &SystemError{}) {
					return err
				}
				errorsSlice = append(errorsSlice, ValidationError{
					Field: propType.Name,
					Err:   err,
				})
			}
		case reflect.Int:
			errorsStack, err := intValidate(int(propValue.Int()), rules)
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
		//nolint:exhaustive
		case reflect.Slice:
			switch propValue.Type().Elem().Kind() {
			case reflect.String:
				for _, val := range propValue.Interface().([]string) {
					err := stringValidate(val, propTagValue)
					if err != nil {
						if errors.As(err, &SystemError{}) {
							return err
						}
						errorsSlice = append(errorsSlice, ValidationError{
							Field: propType.Name,
							Err:   err,
						})
					}
				}
			case reflect.Int:
				errorsStack, err := intValidate(propValue.Interface().([]int), rules)
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
			default:
				return fmt.Errorf("%w: %v", ErrSysUnsupportedSlice, propValue.Type().Elem().Kind())
			}
		default:
			return fmt.Errorf("%w: %v", ErrSysUnsupportedType, propValue.Kind())
		}
	}

	if len(errorsSlice) > 0 {
		return errorsSlice
	}
	return nil
}

func intValidate(value interface{}, rules []Rule) ([]error, error) {
	var errorsSlice []error
	var values []int

	switch v := value.(type) {
	case int:
		values = []int{v}
	case []int:
		values = v
	default:
		return nil, fmt.Errorf("%w: %T", ErrSysUnsupportedType, value)
	}

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

func stringValidate(v string, tag string) error {
	for _, rawRule := range strings.Split(tag, "|") {
		rule := strings.Split(rawRule, ":")

		if len(rule) != 2 {
			return fmt.Errorf("%w: %s", ErrSysInvalidRule, rawRule)
		}

		switch rule[0] {
		case LenRuleTag:
			lenString, err := strconv.Atoi(rule[1])
			if err != nil {
				return fmt.Errorf("%w: %w", ErrSysCantConvertLenValue, err)
			}
			if len(v) != lenString {
				return fmt.Errorf("%w: expected %d, got %d", ErrStringLengthMismatch, lenString, len(v))
			}
		case RegexpRuleTag:
			compiledRegexp, err := getCompiledRegexp(rule[1])
			if err != nil {
				return ErrSysRegexpCompile
			}
			if !compiledRegexp.MatchString(v) {
				return fmt.Errorf("%w: %s", ErrRegexpMatchFailed, rule[1])
			}
		case InRuleTag:
			isMatched := false
			for _, item := range strings.Split(rule[1], ",") {
				if item == v {
					isMatched = true
				}
			}
			if !isMatched {
				return fmt.Errorf("%w: %s in %s", ErrValueNotInList, v, rule[1])
			}
		}
	}
	return nil
}

func getCompiledRegexp(pattern string) (*regexp.Regexp, error) {
	regexpMutex.Lock()
	defer regexpMutex.Unlock()

	if compiled, exists := compiledRegexps[pattern]; exists {
		return compiled, nil
	}

	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	compiledRegexps[pattern] = compiled
	return compiled, nil
}
