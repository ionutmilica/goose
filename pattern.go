package goose

import (
	"fmt"
	"regexp"
)

const (
	STATIC_PATTERN = iota
	REGEX_PATTERN
	PARAM_PATTERN
	WILDCARD_PATTERN
)

type Pattern struct {
	raw       string
	compiled  string
	regex     regexp.Regexp
	kind      int8
	wildcards []string
}

func NewPattern(pattern string) *Pattern {
	patternObj := &Pattern{}
	patternObj.wildcards = make([]string, 0)
	patternObj.raw = pattern
	patternObj.kind = STATIC_PATTERN

	// Prepare the pattern

	patternLen := len(pattern)

	if pattern[0] == '{' && pattern[patternLen-1] == '}' {
		if patternLen < 3 {
			panic(fmt.Sprintf("`%s` pattern is not valid!", pattern))
		}
		var wildcard string
		if isOptionalPattern(pattern) {
			wildcard = pattern[1 : patternLen-2]
		} else {
			wildcard = pattern[1 : patternLen-1]
		}

		patternObj.kind = PARAM_PATTERN
		patternObj.wildcards = append(patternObj.wildcards, wildcard)
		return patternObj
	}

	return patternObj
}

func (self *Pattern) match(against string) (bool, Params) {
	switch self.kind {
	case STATIC_PATTERN:
		if against == self.raw {
			return true, nil
		}
	case PARAM_PATTERN:
		return true, map[string]string{self.wildcards[0]: against}
	}

	return false, nil
}

func appendMap(dst, src map[string]string) map[string]string {
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func isOptionalPattern(pattern string) bool {
	if len(pattern) > 2 && pattern[len(pattern)-2] == '?' {
		return true
	}
	return false
}
