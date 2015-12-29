package goose

import (
	"fmt"
	"strings"
)

type Trie struct {
	root *Node
}

func NewTrie() *Trie {
	return &Trie{NewNode("/")}
}

// Trie pattern inseration
// Returns the leaf node
func (self *Trie) Insert(pattern string, handler Handler) *Node {
	currentNode := self.root

	for _, segment := range strings.Split(pattern, "/") {
		if segment == "" {
			continue
		}
		// Current segment was already added in the tree, so we pass
		if node := currentNode.isChildren(segment); node != nil {
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
		// If the parent of the optional node already has a handler, then it's syntax error
		if currentNode.parent.hasHandler {
			panic(fmt.Sprintf("`%s` node already has a handler and can't be combined with an optional segment!", self))
		}
		currentNode.parent.setHandler(handler)
	}
	currentNode.setHandler(handler)

	return currentNode
}

// Trie lookup
// It returns a node address and url matched params or nil and nil
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
				break
			}
		}
	}

	if currentNode.hasHandler {
		return currentNode, params
	}

	return nil, nil
}

// Dump is used for graph visualization in dot format
func (self *Trie) Dump() {
	groups := map[string]string{}
	walk(self.root, groups)

	format := "graph Router \n{\n%s\n}"

	links := ""

	for parent, child := range groups {
		links = links + fmt.Sprintf("\t\"%s\" -- \"%s\";\n", child, parent)
	}

	fmt.Println(fmt.Sprintf(format, links))
}

// walk is used to map all the node children to their parents in a map[string]string
// It's used to make `dot format` dump of the tree
func walk(node *Node, groups map[string]string) {
	if node == nil {
		return
	}
	for _, child := range node.children {
		walk(child, groups)
		groups[child.String()] = node.String()
	}
}
