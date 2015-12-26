package goose

import (
	"fmt"
	"net/http"
)

var METHODS = map[string]bool{
	"GET":     true,
	"POST":    true,
	"PUT":     true,
	"PATCH":   true,
	"HEAD":    true,
	"OPTIONS": true,
	"DELETE":  true,
}

type Handler func(res http.ResponseWriter, req *http.Request, params Params)

type Group struct {
	pattern string
	handler Handler
}

type Router struct {
	groups []Group
	trees  map[string]*Trie
}

func NewRouter() *Router {
	return &Router{
		groups: make([]Group, 0),
		trees:  make(map[string]*Trie),
	}
}

func (self *Router) Get(pattern string, handler Handler) {
	self.register("GET", pattern, handler)
}

func (self *Router) Post(pattern string, handler Handler) {
	self.register("POST", pattern, handler)
}

func (self *Router) Put(pattern string, handler Handler) {
	self.register("PUT", pattern, handler)
}

func (self *Router) Patch(pattern string, handler Handler) {
	self.register("PATCH", pattern, handler)
}

func (self *Router) Delete(pattern string, handler Handler) {
	self.register("DELETE", pattern, handler)
}

func (self *Router) Any(pattern string, handler Handler) {
	self.register("*", pattern, handler)
}

func (self *Router) register(method, pattern string, handler Handler) {
	if _, ok := METHODS[method]; !ok {
		panic(fmt.Sprintf("Method %s is not valid!", method))
	}

	if _, ok := self.trees[method]; !ok {
		self.trees[method] = NewTrie()
	}

	tree := self.trees[method]
	tree.Insert(pattern, handler)
}

func (self *Router) Group(pattern string, cb func(*Router)) {
	self.groups = append(self.groups, Group{pattern, nil})
	cb(self)
	self.groups = self.groups[:len(self.groups)-1]
}

// Just for some tests
func (self *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if root := self.trees[req.Method]; root != nil {
		path := req.URL.Path

		if node, params := root.Search(path); node != nil {
			node.handler(res, req, params)
			return
		}
	}
	fmt.Fprintln(res, "<h1>Page not found!</h2>")
}
