/**
 * The MIT License (MIT)
 *
 * Copyright (c) 2015 Christian Speckner
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

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
