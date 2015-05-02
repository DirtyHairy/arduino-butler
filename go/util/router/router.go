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
