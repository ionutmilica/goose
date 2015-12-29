package goose

type Params map[string]string

type Node struct {
	parent     *Node
	pattern    *Pattern
	handler    Handler
	hasHandler bool
	children   []*Node
}

// Creates a new node
func NewNode(pattern string) *Node {
	node := &Node{}
	node.hasHandler = false
	node.children = make([]*Node, 0)
	node.pattern = NewPattern(pattern)

	return node
}

// setHandler is used to mark the current node as terminal
func (self *Node) setHandler(handler Handler) {
	self.handler = handler
	self.hasHandler = true
}

// isChildren verifies if a specific pattern or url segment is already present
// as a child node for the current node
func (self Node) isChildren(pattern string) *Node {
	for _, child := range self.children {
		if child.pattern.raw == pattern {
			return child
		}
	}
	return nil
}

// insertChildren will insert a node as a children in a sorted slice.
// It will keep the children ordered by type from lower to higher, eg: 0 - static pattern
func (self *Node) insertChildren(node *Node) {
	i := 0
	for ; i < len(self.children); i++ {
		if node.pattern.kind < self.children[i].pattern.kind {
			break
		}
	}

	if i == len(self.children) {
		self.children = append(self.children, node)
	} else {
		self.children = append(self.children, nil)
		copy(self.children[i+1:], self.children[i:])
		self.children[i] = node
	}
}

// String provides a friendly name for the node
func (self *Node) String() string {
	return self.pattern.raw
}
