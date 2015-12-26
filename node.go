package goose

import (
	"fmt"
	"strings"
)

type Params map[string]string

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
		// Add handler to the node parent is the node is optional
		if newNode.pattern.kind == PARAM_PATTERN && isOptionalPattern(segments[0]) {
			if self.hasHandler {
				panic(fmt.Sprintf("`%s` node already has a handler and can't be combined with an optiona segment!", self))
			}
			self.Insert("", handler)
		}
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
