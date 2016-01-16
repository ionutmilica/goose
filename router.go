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

type Route struct {
	router *Router
	node   *Node
}

func (self *Route) Name(name string) *Route {
	if _, ok := self.router.nameRoutes[name]; ok {
		panic(fmt.Sprintf("Route named `%s` already exists!", name))
	}

	self.router.nameRoutes[name] = self

	return self
}

func (self *Route) Midleware(name string) *Route {
	// Not implemented
	return self
}

type Router struct {
	groups     []Group
	trees      map[string]*Trie
	nameRoutes map[string]*Route
}

// Creates the router struct
func NewRouter() *Router {
	return &Router{
		groups:     make([]Group, 0),
		trees:      make(map[string]*Trie),
		nameRoutes: make(map[string]*Route),
	}
}

// Creates a route for GET request
func (self *Router) Get(pattern string, handler Handler) *Route {
	return self.Handle("GET", pattern, handler)
}

func (self *Router) Head(pattern string, handler Handler) *Route {
	return self.Handle("HEAD", pattern, handler)
}

// Creates a route for POST request
func (self *Router) Post(pattern string, handler Handler) *Route {
	return self.Handle("POST", pattern, handler)
}

func (self *Router) Put(pattern string, handler Handler) *Route {
	return self.Handle("PUT", pattern, handler)
}

func (self *Router) Patch(pattern string, handler Handler) *Route {
	return self.Handle("PATCH", pattern, handler)
}

func (self *Router) Delete(pattern string, handler Handler) *Route {
	return self.Handle("DELETE", pattern, handler)
}

func (self *Router) Options(pattern string, handler Handler) *Route {
	return self.Handle("OPTIONS", pattern, handler)
}

func (self *Router) Any(pattern string, handler Handler) *Route {
	return self.Handle("*", pattern, handler)
}

// Creates a route for a specific request method
func (self *Router) Handle(method, pattern string, handler Handler) *Route {
	methods := make(map[string]bool)

	if method == "*" {
		for m, _ := range METHODS {
			methods[m] = true
		}
	} else {
		methods[method] = true
	}

	if _, ok := METHODS[method]; !ok {
		panic("Unknown Http Method: " + method)
	}

	for m, _ := range methods {
		self.alloc(m)
	}

	return &Route{self, self.trees[method].Insert(pattern, handler)}
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

// Private methods

func (self *Router) alloc(method string) {
	if _, ok := self.trees[method]; !ok {
		self.trees[method] = NewTrie()
	}
}
