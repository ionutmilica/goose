package goose

import (
	"fmt"
	"regexp"
	"strings"
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
	regex     *regexp.Regexp
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

	// Match {something} patterns only for PARAM_PATTERN
	if regexp.MustCompile("^\\{[^\\{\\}:]+\\}$").MatchString(pattern) {
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

	patternObj.compilePattern(pattern)

	return patternObj
}

//Fugly method for compiling params groups into regex
//@todo: Optimize for memory and speed, also deal with unicode characters
func (self *Pattern) compilePattern(pattern string) {
	wildcards := make(map[string]string)
	i := 0
	inWildcard := false
	start := 0

	for i < len(pattern) {
		if inWildcard {
			if pattern[i] == '}' {
				inWildcard = false
				wildcard := pattern[start+1 : i]

				parts := strings.Split(wildcard, ":")

				if len(parts) > 2 {
					panic(fmt.Sprintf("`%s` contains more than 1 modifier (:) in `%s` pattern.", wildcard, pattern))
				}

				if len(parts) > 1 {
					wildcards["{"+wildcard+"}"] = "(?P<" + parts[0] + ">" + parts[1] + ")"
				} else {
					wildcards["{"+wildcard+"}"] = "(?P<" + parts[0] + ">.+)"
				}

				i++
				continue
			}
		} else {
			if pattern[i] == '{' {
				inWildcard = true
				start = i
				i++
				continue
			}
			if pattern[i] == '}' {
				panic(fmt.Sprintf("`%s` pattern is invalid and cannot be compiled!", pattern))
			}
		}
		i++
	}
	for _old, _new := range wildcards {
		pattern = strings.Replace(pattern, _old, _new, -1)
	}

	if self.raw != pattern {
		self.kind = REGEX_PATTERN
		self.compiled = "^" + pattern + "$"
		self.regex = regexp.MustCompile(self.compiled)
	}
}

func (self *Pattern) match(against string) (bool, Params) {
	switch self.kind {
	case STATIC_PATTERN:
		if against == self.raw {
			return true, nil
		}
	case REGEX_PATTERN:
		parts := self.regex.FindStringSubmatch(against)
		if len(parts) > 0 {
			params := make(map[string]string, len(parts)-1)
			for i, name := range self.regex.SubexpNames() {
				if i == 0 {
					continue
				}
				params[name] = parts[i]
			}
			return true, params
		}
		fmt.Println(parts, self.regex.SubexpNames())
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
