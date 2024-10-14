package hw09structvalidator

import (
	"fmt"
	"regexp"
	"sync"
)

var (
	compiledRegexps = make(map[string]*regexp.Regexp)
	regexpMutex     = sync.RWMutex{}
)

type RegexpRule struct {
	Value string
}

func (r RegexpRule) Validate(v interface{}) error {
	v, ok := v.(string)
	if !ok {
		return fmt.Errorf("%w: %T", ErrSysUnsupportedType, v)
	}

	compiledRegexp, err := getCompiledRegexp(r.Value)
	if err != nil {
		return ErrSysRegexpCompile
	}
	if !compiledRegexp.MatchString(v.(string)) {
		return fmt.Errorf("%w: %s", ErrRegexpMatchFailed, r.Value)
	}
	return nil
}

func getCompiledRegexp(pattern string) (*regexp.Regexp, error) {
	regexpMutex.RLock()
	if compiled, exists := compiledRegexps[pattern]; exists {
		regexpMutex.RUnlock()
		return compiled, nil
	}
	regexpMutex.RUnlock()

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
