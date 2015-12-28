package goose

import (
	"fmt"
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
	return self.root.Insert(pattern, handler)
}

func (self *Trie) Search(pattern string) (*Node, Params) {
	params := make(Params, 0)
	segments := splitIntoSegments(pattern)
	currentNode := self.root
	numSegments := len(segments)
	i := 0

	for i < numSegments {
		for _, child := range currentNode.children {
			if child.pattern.match(segments[i], params) {
				currentNode = child
			}
		}
		i++
	}

	if currentNode.hasHandler {
		return currentNode, params
	}

	return nil, nil
}

func (self *Trie) Dump() {
	groups := map[string]string{}
	walk(self.root, groups)

	fmt.Println(self.root.children)

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
