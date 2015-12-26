package goose

type Trie struct {
	root *Node
}

func NewTrie() *Trie {
	return &Trie{}
}

func (self *Trie) Add(pattern string, handler Handler) {
	if self.root == nil {
		self.root = NewNode()
	}
	self.root.Insert(pattern, handler)
}
