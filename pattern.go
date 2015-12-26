package goose

import (
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

	// Prepare the pattern

	if pattern[0] == '{' && pattern[len(pattern)-1] == '}' {
		patternObj.kind = PARAM_PATTERN
		patternObj.wildcards = append(patternObj.wildcards, pattern[1:len(pattern)-1])
		return patternObj
	}

	patternObj.kind = STATIC_PATTERN

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
