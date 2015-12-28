package goose

import (
	"fmt"
	"strings"
)

type Trie struct {
	root *Node
}

func NewTrie() *Trie {
	return &Trie{}
}

func (self *Trie) Insert(pattern string, handler Handler) *Node {
	if self.root == nil {
		self.root = NewNode("/")
	}

	currentNode := self.root

	for _, segment := range strings.Split(pattern, "/") {
		if segment == "" {
			continue
		}
		// Current segment was already added in the tree, so we pass
		if node := currentNode.inChildren(segment); node != nil {
			currentNode = node
		} else {
			// We should create the new node and insert it by priority
			newNode := NewNode(segment)
			newNode.parent = currentNode
			currentNode.insertChildren(newNode)
			currentNode = newNode
		}
	}

	if currentNode.pattern.kind == PARAM_PATTERN && currentNode.pattern.isOptional {
		if currentNode.parent.hasHandler {
			panic(fmt.Sprintf("`%s` node already has a handler and can't be combined with an optiona segment!", self))
		}
		currentNode.parent.setHandler(handler)
	}
	currentNode.setHandler(handler)

	return currentNode
}

func (self *Trie) Search(pattern string) (*Node, Params) {
	params := make(Params, 0)
	currentNode := self.root

	for _, segment := range strings.Split(pattern, "/") {
		if segment == "" {
			continue
		}
		for _, child := range currentNode.children {
			if child.pattern.match(segment, params) {
				currentNode = child
			}
		}
	}

	if currentNode.hasHandler {
		return currentNode, params
	}

	return nil, nil
}

func (self *Trie) Dump() {
	groups := map[string]string{}
	walk(self.root, groups)

	format := `graph Router {
%s
}`

	links := ""

	for parent, child := range groups {
		links = links + fmt.Sprintf("   \"%s\" -- \"%s\";\n", child, parent)
	}

	fmt.Println(fmt.Sprintf(format, links))
}

func walk(node *Node, groups map[string]string) {
	if node == nil {
		return
	}
	for _, child := range node.children {
		walk(child, groups)
		groups[child.String()] = node.String()
	}
}
