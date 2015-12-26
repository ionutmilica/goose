package goose

import (
	"fmt"
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

type Handler interface{}

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
	tree.Add(pattern, handler)
}

func (self *Router) Group(pattern string, cb func(*Router)) {
	self.groups = append(self.groups, Group{pattern, cb})
	cb(self)
	self.groups = self.groups[:len(self.groups)-1]
}
