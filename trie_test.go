package goose

import (
	"testing"
)

func makeTrie(routes []string) *Trie {
	trie := NewTrie()
	for _, route := range routes {
		trie.Insert(route, nil)
	}
	return trie
}

// Check if two maps have the same elements
// Elements order is neglected
func equalMaps(a, b map[string]string) bool {
	for k, v := range a {
		if v2, ok := b[k]; !ok || v != v2 {
			return false
		}
	}
	for k, v := range b {
		if v2, ok := a[k]; !ok || v != v2 {
			return false
		}
	}
	return true
}

// Having some routes and expected results this function will make all the tests
func routesWithParamsTest(t *testing.T, routes []string, toMatch map[string]map[string]string) {
	trie := makeTrie(routes)

	for match, expect := range toMatch {
		var node *Node
		var params Params

		node, params = trie.Search(match)

		if node == nil && expect != nil {
			t.Errorf("Tried to match %s but failed!", match)
		}

		if expect != nil && !equalMaps(params, expect) {
			t.Errorf("For [%s], expected params: [%s] but got [%s]", match, expect, params)
		}
	}
}

func TestStaticRoutes(t *testing.T) {
	routes := []string{
		"post/hello-world",
		"post/another-post",
		"post/more-post",
		"pages/about",
		"pages/contact-us",
		"/",
		"about-us",
		"files/index.html",
	}

	trie := makeTrie(routes)
	for _, route := range routes {
		if node, _ := trie.Search(route); node == nil {
			t.Errorf("Route %s not matched!", route)
		}
	}
}

func TestParamRoutes(t *testing.T) {
	routes := []string{
		"post/hello-world",
		"post/{post}",
		"pages/{page}",
	}

	toMatch := map[string]map[string]string{
		"post/hello-world":     map[string]string{},
		"post/my-awesome-post": map[string]string{"post": "my-awesome-post"},
		"posts/some-post":      nil,
		"post/アニメ":             map[string]string{"post": "アニメ"},
		"post/holla/ups":       nil,
		"post/something/10":    nil,
	}

	routesWithParamsTest(t, routes, toMatch)
}

func TestMoreParamRoutes(t *testing.T) {
	routes := []string{
		"users/{id}",
		"users/profile/{user}",
		"call/{me}/{baby}",
	}

	toMatch := map[string]map[string]string{
		"/users/ionut": map[string]string{"id": "ionut"},
		"/users/10":    map[string]string{"id": "10"},
		// Take in consideration this match
		"/users/profile":       nil,
		"/users/profile/ionut": map[string]string{"user": "ionut"},
		"/users/profile/10":    map[string]string{"user": "10"},
		"/users/profile/アニメ":   map[string]string{"user": "アニメ"},
		"/call":                nil,
		"/call/me":             nil,
		"/call/me/baby/":       map[string]string{"me": "me", "baby": "baby"},
	}

	routesWithParamsTest(t, routes, toMatch)
}

func TestOptionalParamRoutes(t *testing.T) {
	routes := []string{
		"users/{id?}",
		"users/profile/{user}",
	}

	toMatch := map[string]map[string]string{
		"users/ionut":         map[string]string{"id": "ionut"},
		"users/10":            map[string]string{"id": "10"},
		"users":               map[string]string{},
		"users/profile/ionut": map[string]string{"user": "ionut"},
	}

	routesWithParamsTest(t, routes, toMatch)
}

func TestRegexParamRoutes(t *testing.T) {
	routes := []string{
		"post/{title}-{id}",
		"{something}-us",
	}

	toMatch := map[string]map[string]string{
		"post/hello-summer-10": map[string]string{"title": "hello-summer", "id": "10"},
		"about-us":             map[string]string{"something": "about"},
	}

	routesWithParamsTest(t, routes, toMatch)
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
