package goose

import (
	"net/http"
	"testing"
)

func TestNodeInitialization(t *testing.T) {
	node := NewNode("/")

	if node == nil {
		t.Error("Node `/` was not created!")
	}

	if node.pattern == nil {
		t.Error("Node `/` got not pattern object.")
	}

	if node.children == nil {
		t.Error("Node `/` children property was not initialized.")
	}
}

func TestString(t *testing.T) {
	root := NewNode("/")

	if root.pattern.raw != root.String() {
		t.Error("root.String() is not equal with root.pattern.raw")
	}
}

func TestSetHandler(t *testing.T) {
	node := NewNode("/")
	node.setHandler(func(res http.ResponseWriter, req *http.Request, params Params) {

	})
	if node.handler == nil {
		t.Error("Node `/` has no handler attached!")
	}

	if node.hasHandler == false {
		t.Error("Node `/` hasHandler property is false!")
	}
}

func TestInsertChildren(t *testing.T) {
	root := NewNode("/")
	id := NewNode("{id}")
	profile := NewNode("profile")

	root.insertChildren(id)
	root.insertChildren(profile)

	if root.children[0] != profile {
		t.Errorf("%s node should be children 0 of root!", profile)
	}
}
