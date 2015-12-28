package goose

import (
	"fmt"
)

type Params map[string]string

type Node struct {
	parent     *Node
	pattern    *Pattern
	handler    Handler
	hasHandler bool
	children   []*Node
}

func NewNode(pattern string) *Node {
	node := &Node{}
	node.hasHandler = false
	node.children = make([]*Node, 0)
	node.pattern = NewPattern(pattern)

	return node
}

func (self *Node) String() string {
	return self.pattern.raw
}

func (self *Node) setHandler(handler Handler) {
	self.handler = handler
	self.hasHandler = true
}

// Naive implementation without optimization toughts
// Will be revised when when is completed
func (self *Node) insert(segments []string, handler Handler) *Node {
	if len(segments) == 0 {
		self.setHandler(handler)
		return self
	}

	var newNode *Node
	for _, n := range self.children {
		if n.pattern.raw == segments[0] {
			//if n.pattern.match(segments[0], make(Params, 0)) {
			newNode = n
			break
		}
	}

	if newNode == nil {
		newNode = NewNode(segments[0])
		newNode.parent = self
		self.children = append(self.children, newNode)
	}

	if len(segments) == 1 {
		// Add handler to the parrent node if the current node is optional
		if newNode.pattern.kind == PARAM_PATTERN && isOptionalPattern(segments[0]) {
			if self.hasHandler {
				panic(fmt.Sprintf("`%s` node already has a handler and can't be combined with an optiona segment!", self))
			}
			self.setHandler(handler)
		}
		newNode.setHandler(handler)

		return newNode
	}

	return newNode.insert(segments[1:], handler)
}

func (self *Node) Insert(pattern string, handler Handler) *Node {
	return self.insert(splitIntoSegments(pattern), handler)
}
