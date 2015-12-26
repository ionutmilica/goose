package goose

import (
	"testing"
)

func TestStaticPattern(t *testing.T) {
	tests := []string{
		"user", "test", "1005", "пиво",
	}

	for _, tt := range tests {
		if pattern := NewPattern(tt); pattern.kind != STATIC_PATTERN {
			t.Errorf("%s was not compiled as static pattern!", tt)
		}
	}
}

// This tests patterns that contain only the param, ex:  {user}
func TestOnlyParamPattern(t *testing.T) {
	tests := []string{
		"{user}", "{id}",
	}

	for _, tt := range tests {
		if pattern := NewPattern(tt); pattern.kind != PARAM_PATTERN {
			t.Errorf("%s was not compiled as param pattern!", tt)
			if len(pattern.wildcards) == 0 || pattern.wildcards[0] != tt[1:len(tt)-1] {
				t.Errorf("Wildcard `%s` for pattern `%s` not found or not matching!", tt[1:len(tt)], tt)
			}
		}
	}
}
