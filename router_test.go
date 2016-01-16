package goose

import (
	"testing"
)

func TestGooseInitialisation(t *testing.T) {
	router := NewRouter()

	if router == nil {
		t.Error("Goose was not instantiated!")
	}

	if router.groups == nil {
		t.Error("Router groups property was not instantiated!")
	}
}

func TestRouteRegister(t *testing.T) {
	router := NewRouter()
	if route := router.Get("/", nil); route == nil {
		t.Error("Router.Get returned nil")
	}
	if route := router.Head("/", nil); route == nil {
		t.Error("Router.Head returned nil")
	}
	if route := router.Post("/", nil); route == nil {
		t.Error("Router.Post returned nil")
	}
	if route := router.Put("/", nil); route == nil {
		t.Error("Router.Put returned nil")
	}
	if route := router.Patch("/", nil); route == nil {
		t.Error("Router.Patch returned nil")
	}
	if route := router.Delete("/", nil); route == nil {
		t.Error("Router.Delete returned nil")
	}
	if route := router.Options("/", nil); route == nil {
		t.Error("Router.Options returned nil")
	}
}
