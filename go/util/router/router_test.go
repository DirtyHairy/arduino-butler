package router

import (
	"github.com/DirtyHairy/arduino-butler/go/util/mock"
	"net/http"
	"net/url"
	"testing"
)

type stateT struct {
	handler1Called int
	matches        []string
	handler2Called int
}

var state stateT

func (state *stateT) reset() {
	state.handler1Called = 0
	state.handler2Called = 0
	state.matches = nil
}

func prepareRouter(router *Router) {
	router.AddRoute("^/foo/bar/(\\d+)$", HandlerFunction(func(response http.ResponseWriter, request *http.Request, matches []string) {
		state.handler1Called++
		state.matches = matches
	}))

	router.AddRoute("^/bar/foo/(\\d+)$", HandlerFunction(func(response http.ResponseWriter, request *http.Request, matches []string) {
		state.handler2Called++
		state.matches = matches
	}))
}

func TestRoute1(t *testing.T) {
	router := CreateRouter(10)

	prepareRouter(&router)
	state.reset()

	var request http.Request
	request.URL, _ = url.Parse("/foo/bar/1")
	response := mock.CreateMockResponseWriter()

	router.ServeHTTP(&response, &request)

	if state.handler1Called != 1 {
		t.Error("Handler 1 should have been called")
	}

	if state.handler2Called != 0 {
		t.Error("Handler 2 should not have been called")
	}

	if state.matches == nil || state.matches[1] != "1" {
		t.Error("Matching failed")
	}
}

func TestRoute2(t *testing.T) {
	router := CreateRouter(10)

	prepareRouter(&router)
	state.reset()

	var request http.Request
	request.URL, _ = url.Parse("/bar/foo/2")
	response := mock.CreateMockResponseWriter()

	router.ServeHTTP(&response, &request)

	if state.handler1Called != 0 {
		t.Error("Handler 1 should not have been called")
	}

	if state.handler2Called != 1 {
		t.Error("Handler 2 should have been called")
	}

	if state.matches == nil || state.matches[1] != "2" {
		t.Error("Matching failed")
	}
}

func TestInvalidRoute(t *testing.T) {
	router := CreateRouter(10)

	prepareRouter(&router)
	state.reset()

	var request http.Request
	request.URL, _ = url.Parse("/bar/foo/")
	response := mock.CreateMockResponseWriter()

	router.ServeHTTP(&response, &request)

	if state.handler1Called != 0 {
		t.Error("Handler 2 should not have been called")
	}

	if state.handler2Called != 0 {
		t.Error("Handler 1 should not have been called")
	}

	if response.Code != http.StatusNotFound {
		t.Errorf("Result should have been 404 - not found, got %d instead", response.Code)
	}
}
