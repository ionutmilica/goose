package goose

import (
	"testing"
)

func TestStaticPath(t *testing.T) {
	tree := NewNode()

	tree.Insert("", "root")
	tree.Insert("/users/profile/info", "/users/profile/info")
	tree.Insert("/users/profile", "users/profile")
	tree.Insert("/users/test", "users/test")
	tree.Insert("/user", "user")
	if node, _ := tree.Search(""); node == nil {
		t.Error("Cannot find root node with empty space")
	}

	if node, _ := tree.Search("/"); node == nil {
		t.Error("Cannot find root node with / pattern")
	}

	if node, _ := tree.Search("/users/profile/info"); node == nil {
		t.Error("No node was found with /users/profile/info pattern")
	}

	if node, _ := tree.Search("user"); node == nil {
		t.Error("No node was found with user pattern")
	}

	if node, _ := tree.Search("/users/test"); node == nil {
		t.Error("No node was found with users/test pattern")
	}

	if node, _ := tree.Search("/users/profile"); node == nil {
		t.Error("No node was found with users/profile pattern")
	}
}

func TestParams(t *testing.T) {
	tree := NewNode()
	tree.Insert("users/{id}", "user_id")
	tree.Insert("users/profile/{user}", "test_dsa")
	tree.Insert("some/{more}/{even}", "some_more_even")

	node, params := tree.Search("users/10")
	if node == nil {
		t.Error("users/10 was not matched")
	}

	node, params = tree.Search("users/profile/ionut")
	if node == nil {
		t.Error("users/profile/{user} was not matched")
	}
	if _, ok := params["user"]; !ok {
		t.Error("{user} not found in match for users/profile/{user}")
	}

	node, params = tree.Search("some/ionut/milica")
	if node == nil {
		t.Error("test/adsa was not matched")
	}
	if _, ok := params["more"]; !ok {
		t.Error("{more} not found in match for some/{more}/{even}")
	}
	if _, ok := params["even"]; !ok {
		t.Error("{even} not found in match for some/{more}/{even}")
	}

}

func TestOptionalParameter(t *testing.T) {
	tree := NewNode()
	tree.Insert("users/{id?}", "user_id")
	tree.Insert("users/profile/{user}", "test_dsa")

	node, params := tree.Search("users")

	if node == nil {
		t.Error("No results found for `users` although we have an optional parameter")
	}

	node, params = tree.Search("users/10")

	if node == nil {
		t.Error("No results found for `users/10`")
	}

	if _, ok := params["id"]; !ok {
		t.Error("{id} not populated despite we had users/10 and users/{id?} pattern")
	}
}
