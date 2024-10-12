package hw09structvalidator

import (
	"fmt"
	"strconv"
	"strings"
)

type Rule interface {
	Validate(interface{}) error
}

func parseRules(tag string) ([]Rule, error) {
	var rules []Rule
	for _, rawRule := range strings.Split(tag, "|") {
		rule := strings.Split(rawRule, ":")

		if len(rule) != 2 {
			return rules, fmt.Errorf("%w: %s", ErrSysInvalidRule, tag)
		}

		switch rule[0] {
		case MinRuleTag:
			minValue, err := strconv.Atoi(rule[1])
			if err != nil {
				return rules, fmt.Errorf("%w: %w", ErrSysCantConvertMinValue, err)
			}
			rules = append(rules, MinRule{Value: minValue})
		case MaxRuleTag:
			maxValue, err := strconv.Atoi(rule[1])
			if err != nil {
				return rules, fmt.Errorf("%w: %w", ErrSysCantConvertMaxValue, err)
			}
			rules = append(rules, MaxRule{Value: maxValue})

		case InRuleTag:
			enumValues := strings.Split(rule[1], ",")
			if len(enumValues) == 0 {
				return rules, fmt.Errorf("%w: %s", ErrSysInvalidRule, tag)
			}
			rules = append(rules, InRule{Values: enumValues})
		case LenRuleTag:
			lenValue, err := strconv.Atoi(rule[1])
			if err != nil {
				return rules, fmt.Errorf("%w: %w", ErrSysCantConvertLenValue, err)
			}
			rules = append(rules, LenRule{Value: lenValue})
		case RegexpRuleTag:
			rules = append(rules, RegexpRule{Value: rule[1]})
		}
	}
	return rules, nil
}
