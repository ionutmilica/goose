package goose

type Trie struct {
	root *Node
}

func NewTrie() *Trie {
	return &Trie{}
}

func (self *Trie) Insert(pattern string, handler Handler) *Node {
	if self.root == nil {
		self.root = NewNode()
	}
	return self.root.Insert(pattern, handler)
}

func (self *Trie) Search(pattern string) (*Node, Params) {
	return self.root.Search(pattern)
}
