package goose

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

func (self Node) inChildren(pattern string) *Node {
	for _, child := range self.children {
		if child.pattern.raw == pattern {
			return child
		}
	}
	return nil
}

func (self *Node) insertChildren(node *Node) {
	i := 0
	for ; i < len(self.children); i++ {
		if node.pattern.kind < self.children[i].pattern.kind {
			break
		}
	}
	if i == len(node.children) {
		self.children = append(self.children, node)
	} else {
		self.children = append(self.children[:i], append([]*Node{node}, self.children[i:]...)...)
	}
}
