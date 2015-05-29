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
	"github.com/DirtyHairy/arduino-butler/go/util"
	"net/http"
	"regexp"
)

type Handler interface {
	Handle(response http.ResponseWriter, request *http.Request, matches []string)
}

type HandlerFunction func(http.ResponseWriter, *http.Request, []string)

type route struct {
	expression string
	regex      *regexp.Regexp
	handler    Handler
}

type Router struct {
	routes []route
}

func (router Router) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	var matches []string
	var route route

	for _, route = range router.routes {
		matches = route.regex.FindStringSubmatch(request.URL.Path)

		if matches != nil {
			break
		}
	}

	if matches != nil {
		route.handler.Handle(response, request, matches)
	} else {
		http.NotFound(response, request)
	}
}

func (router *Router) AddRoute(expression string, handler Handler) {
	router.routes = append(router.routes, route{expression, util.CompileRegex(expression), handler})
}

func CreateRouter(initialRouteLimit uint) Router {
	return Router{make([]route, 0, initialRouteLimit)}
}

func (fn HandlerFunction) Handle(response http.ResponseWriter, request *http.Request, matches []string) {
	fn(response, request, matches)
}
