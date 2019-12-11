package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type key string

const (
	pathContextKey       key = "muxer_pathContextKeyFooBar"
	pathOffsetContextKey key = "muxer_pathOffsetContextKeyFooBar"
)

type Muxer struct {
	routes *Routes
}

func NewMuxer() *Muxer {
	routes := NewRouter()
	return &Muxer{&routes}
}

func wrapHandler(path string, h http.Handler) http.Handler {

	// path is /foo/bar/baz
	// offset would then be 0
	//
	// path is /something/another
	// offset would then be 3
	//
	// How do we know? We record the previous path!
	//
	// path is /even/more
	// offset would be 5
	//
	// How do we know? We record the previous path, and we add the offset!

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		previousPath, ok := r.Context().Value(pathContextKey).(string)
		if !ok {
			previousPath = "/"
		}
		pathOffset, ok := r.Context().Value(pathOffsetContextKey).(int)
		if !ok {
			pathOffset = 0
		}

		// The first slash is a distractino.
		previousPathNoLeadingSlash := previousPath[1:]

		components := strings.Split(previousPathNoLeadingSlash, "/")
		components = components[1:]

		newOffset := pathOffset + len(components)

		// Store the route's path (not the request path)
		ctx := context.WithValue(r.Context(), pathContextKey, path)
		ctx = context.WithValue(ctx, pathOffsetContextKey, newOffset)

		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})

}

func Params(r *http.Request) map[string]string {
	requestPathComponents := r.URL.Path

	ctx := r.Context()

	routePath, ok := ctx.Value(pathContextKey).(string)
	if !ok {
		return make(map[string]string)
	}
	pathOffset, ok := ctx.Value(pathOffsetContextKey).(int)
	if !ok {
		return make(map[string]string)
	}

	fmt.Println(routePath)
	fmt.Println(pathOffset)

	routePathComponents := strings.Split(routePath, "/")[1:]
	pathComponents := strings.Split(requestPathComponents, "/")[1+pathOffset:]

	if len(routePathComponents) != len(pathComponents) {
		panic("Something went wrong")
	}

	result := make(map[string]string)
	for i := 0; i < len(routePathComponents); i++ {
		if routePathComponents[i][0] == ':' {
			result[routePathComponents[i][1:]] = pathComponents[i]
		}
	}

	return result
}

func (m *Muxer) addHandlerMethod(path string, method string, h http.Handler) {
	handler, ok := m.routes.Get(path).(*RouteHandler)
	h = wrapHandler(path, h)
	if handler == nil || !ok {
		m.routes.Add(path, &RouteHandler{method: h})
	} else {
		(*handler)[method] = h
	}
}

func (m *Muxer) AddGetHandler(path string, h http.Handler) {
	m.addHandlerMethod(path, "GET", h)
}
func (m *Muxer) AddPostHandler(path string, h http.Handler) {
	m.addHandlerMethod(path, "POST", h)
}
func (m *Muxer) AddPutHandler(path string, h http.Handler) {
	m.addHandlerMethod(path, "PUT", h)
}
func (m *Muxer) AddDeleteHandler(path string, h http.Handler) {
	m.addHandlerMethod(path, "DELETE", h)
}

func (m *Muxer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler, ok := m.routes.Get(req.URL.Path).(*RouteHandler)
	if handler == nil || !ok {
		w.WriteHeader(404)
		w.Write([]byte("Not found"))
	} else {
		// ctx := context.WithValue(req.Context(), pathContextKey, req.URL.Path)
		// req = req.WithContext(ctx)
		handler.ServeHTTP(w, req)
	}
}
