package goose

import (
	"testing"
)

func TestStaticPath(t *testing.T) {
	tree := NewTrie()

	tree.Insert("", nil)
	tree.Insert("/users/profile/info", nil)
	tree.Insert("/users/profile", nil)
	tree.Insert("/users/test", nil)
	tree.Insert("/user", nil)
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
	tree := NewTrie()
	tree.Insert("users/{id}", nil)
	tree.Insert("users/profile/{user}", nil)
	tree.Insert("some/{more}/{even}", nil)

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
	tree := NewTrie()
	tree.Insert("users/{id?}", nil)
	tree.Insert("users/profile/{user}", nil)

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

func TestRegexParameters(t *testing.T) {
	tree := NewTrie()
	tree.Insert("post/{title}-{id}", nil)

	node, params := tree.Search("post/somepost-12312")

	if node == nil {
		t.Error("`post/somepost-12312` was not matched by `post/{title}-{id}` pattern!")
	}

	if val, ok := params["title"]; !ok || val != "somepost" {
		t.Error("Param `title` was not found or not equal with `somepost`!")
	}
	if val, ok := params["id"]; !ok || val != "12312" {
		t.Error("Param `id` was not found or not equal with `12312`!")
	}

}

func BenchmarkTweenty(b *testing.B) {
	route := "/{a}/{b}/{c}/{d}/{e}/{f}/{g}/{h}/{i}/{j}/{k}/{l}/{m}/{n}/{o}/{p}/{q}/{r}/{s}/{t}"

	tree := NewTrie()
	tree.Insert(route, nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if node, _ := tree.Search("/a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t"); node != nil {

		}
	}
}
