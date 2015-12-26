package goose

import (
	"regexp"
	"strings"
)

const (
	STATIC_PATTERN = iota
	REGEX_PATTERN
	PARAM_PATTERN
	WILDCARD_PATTERN
)

func appendMap(dst, src map[string]string) map[string]string {
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

type Params map[string]string

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

	// Prepare the pattern

	if pattern[0] == '{' && pattern[len(pattern)-1] == '}' {
		patternObj.kind = PARAM_PATTERN
		patternObj.wildcards = append(patternObj.wildcards, pattern[1:len(pattern)-1])
	} else {
		patternObj.kind = STATIC_PATTERN
	}

	patternObj.raw = pattern

	//

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

type Node struct {
	parent     *Node
	pattern    *Pattern
	handler    Handler
	hasHandler bool
	children   []*Node
}

func NewNode() *Node {
	node := &Node{}
	node.hasHandler = false
	node.children = make([]*Node, 0)

	return node
}

func (self *Node) String() string {
	return self.pattern.raw
}

// Naive implementation without optimization toughts
// Will be revised when when is completed
func (self *Node) Insert(pattern string, handler Handler) *Node {
	if pattern == "" || pattern == "/" {
		self.handler = handler
		self.hasHandler = true
		return self
	}

	pattern = strings.Trim(pattern, "/")
	segments := strings.Split(pattern, "/")

	// Add node as terminal
	var newNode *Node
	for _, n := range self.children {
		if ok, _ := n.pattern.match(segments[0]); ok {
			newNode = n
		}
	}

	if newNode == nil {
		newNode = NewNode()
		newNode.pattern = NewPattern(segments[0])
		newNode.parent = self
		self.children = append(self.children, newNode)
	}

	if len(segments) == 1 {
		return newNode.Insert("", handler)
	}

	pattern = strings.Join(segments[1:], "/")
	return newNode.Insert(pattern, handler)
}

func (self *Node) Search(pattern string) (*Node, Params) {
	params := make(Params)
	if pattern == "" || pattern == "/" {
		if self.hasHandler {
			return self, params
		}
		return nil, params
	}

	pattern = strings.Trim(pattern, "/")
	segments := strings.Split(pattern, "/")

	numSegments := len(segments)

	for _, node := range self.children {
		if ok, newParams := node.pattern.match(segments[0]); ok {
			params = appendMap(params, newParams)
			if numSegments == 1 && node.hasHandler {
				return node, params
			}
			pattern := strings.Join(segments[1:], "/")
			node, newParams := node.Search(pattern)
			return node, appendMap(params, newParams)
		}
	}

	return nil, params
}
